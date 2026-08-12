# Codebase Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Correct blackjack round behavior and concurrent API state handling, repair the trainer flow and quality checks, and delete confirmed dead code.

**Architecture:** Keep the existing Go engine/API and Svelte client. Make the game engine authoritative for rules, serialize access inside the in-memory session store, and let the client mirror server state while guarding asynchronous hints.

**Tech Stack:** Go 1.25+, `net/http`, Svelte 4, TypeScript 5, Vitest, ESLint, golangci-lint 2.

## Global Constraints

- Six decks; dealer stands on soft 17.
- Double down only on the first two cards, including after a split.
- No surrender or insurance.
- Split 21 is not natural blackjack.
- No new dependencies, persistence, betting features, or visual redesign.
- Every behavior change starts with a failing regression test.

---

### Task 1: Correct engine rules and round lifecycle

**Files:**
- Modify: `internal/hand/hand.go`
- Modify: `internal/game/game.go`
- Modify: `internal/game/game_test.go`

**Interfaces:**
- Consumes: existing `Game.StartRound`, `Game.ApplyMove`, and `Hand.IsBlackjack` behavior.
- Produces: `Hand.FromSplit bool`, `Game.AbandonRound()`, two dealer cards at deal time, legal-double enforcement, and correct split-21 outcomes.

- [ ] **Step 1: Add failing engine regression tests**

Add tests that set a deterministic shoe from inside package `game`:

```go
func TestGame_StartRoundDealsDealerHoleCard(t *testing.T) {
	g := NewGame(player.NewPlayer(), dealer.NewDealer())
	g.StartRound()
	if got := len(g.Dealer.Hand.Cards); got != 2 {
		t.Fatalf("dealer cards = %d, want 2", got)
	}
}

func TestGame_ApplyMoveRejectsDoubleAfterHit(t *testing.T) {
	p := player.NewPlayer()
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Three), card.NewCard(card.Four)})
	d := dealer.NewDealer()
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})
	g := NewGame(p, d)
	g.shoe = shoe.NewShoe([]deck.Deck{{Cards: []card.Card{
		card.NewCard(card.Two), card.NewCard(card.Three),
	}}})

	if err := g.ApplyMove(strategy.Hit); err != nil {
		t.Fatal(err)
	}
	if err := g.ApplyMove(strategy.DoubleDown); !errors.Is(err, ErrInvalidMove) {
		t.Fatalf("double error = %v, want %v", err, ErrInvalidMove)
	}
}

func TestGame_setOutcomesSplitTwentyOneIsWin(t *testing.T) {
	p := player.NewPlayer()
	split21 := hand.NewHand([]card.Card{card.NewCard(card.Ace), card.NewCard(card.Ten)})
	split21.FromSplit = true
	p.ResolvedHands = []*hand.Hand{split21}
	d := dealer.NewDealer()
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Nine)})
	g := NewGame(p, d)

	g.setOutcomes()
	if got := g.Outcomes[0]; got != OutcomeWin {
		t.Fatalf("outcome = %q, want %q", got, OutcomeWin)
	}
}

func TestGame_AbandonRound(t *testing.T) {
	g := NewGame(player.NewPlayer(), dealer.NewDealer())
	g.StartRound()
	g.AbandonRound()
	if g.RoundState != RoundStateComplete || g.Player.ActiveHand != nil {
		t.Fatalf("abandoned game remains active: %#v", g)
	}
}
```

- [ ] **Step 2: Run the focused tests and verify RED**

Run: `go test ./internal/game -run 'StartRoundDealsDealerHoleCard|RejectsDoubleAfterHit|SplitTwentyOneIsWin|AbandonRound'`

Expected: failures because the dealer has one card, double succeeds, `FromSplit` and `AbandonRound` do not exist, and split 21 is classified as blackjack.

- [ ] **Step 3: Implement the minimum engine changes**

Add `FromSplit bool` to `Hand`. In `StartRound`, draw two dealer cards. In `double`, reject any hand whose card count is not two. In `split`, mark both new hands as split hands. In `setOutcomes`, require `!playerHand.FromSplit` for `OutcomeBlackjack`. Add:

```go
func (g *Game) AbandonRound() {
	if g.RoundState != RoundStateActive {
		return
	}
	g.Player.RefreshHand(nil)
	g.Outcomes = nil
	g.RoundState = RoundStateComplete
}
```

- [ ] **Step 4: Verify GREEN and the full engine suite**

Run: `go test ./internal/game ./internal/hand ./internal/strategy/...`

Expected: all tests pass.

- [ ] **Step 5: Commit the engine fix**

```bash
git add internal/hand/hand.go internal/game/game.go internal/game/game_test.go
git commit -m "fix: enforce blackjack round rules"
```

### Task 2: Serialize sessions, expire idle games, and tighten API routing

**Files:**
- Modify: `api/store/session.go`
- Create: `api/store/session_test.go`
- Modify: `api/server.go`
- Modify: `api/handler/game.go`
- Modify: `api/handler/deal.go`
- Modify: `api/handler/hint.go`
- Modify: `api/handler/move.go`
- Create: `api/handler/abandon.go`
- Create: `api/handler/handler_test.go`

**Interfaces:**
- Consumes: `Game.AbandonRound()` from Task 1.
- Produces: `SessionStore.WithGame(id, func(*game.Game)) bool`, 24-hour opportunistic expiry, and `POST /api/game/{id}/abandon`.

- [ ] **Step 1: Write failing session-store tests**

Create `api/store/session_test.go` with a controllable clock and a two-goroutine serialization check:

```go
func TestSessionStoreExpiresIdleGames(t *testing.T) {
	now := time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC)
	store := newSessionStore(func() time.Time { return now })
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
	store.Create(g)
	now = now.Add(24*time.Hour + time.Second)
	if store.WithGame(g.ID, func(*game.Game) {}) {
		t.Fatal("expired game still exists")
	}
}

func TestSessionStoreSerializesOneGame(t *testing.T) {
	store := NewSessionStore()
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
	store.Create(g)
	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	go store.WithGame(g.ID, func(*game.Game) {
		close(entered)
		<-release
	})
	<-entered
	attempted := make(chan struct{})
	go func() {
		close(attempted)
		store.WithGame(g.ID, func(*game.Game) { close(done) })
	}()
	<-attempted
	select {
	case <-done:
		t.Fatal("second operation entered while first held the game")
	default:
	}
	close(release)
	<-done
}
```

- [ ] **Step 2: Run store tests and verify RED**

Run: `go test ./api/store`

Expected: compile failure because `newSessionStore` and `WithGame` do not exist.

- [ ] **Step 3: Implement per-game locking and expiry**

Replace the map of raw games with locked entries:

```go
const sessionTTL = 24 * time.Hour

type session struct {
	mu         sync.Mutex
	game       *game.Game
	lastAccess time.Time
}

type SessionStore struct {
	mu    sync.Mutex
	games map[string]*session
	now   func() time.Time
}

func NewSessionStore() *SessionStore { return newSessionStore(time.Now) }

func newSessionStore(now func() time.Time) *SessionStore {
	return &SessionStore{games: make(map[string]*session), now: now}
}

func (s *SessionStore) WithGame(id string, fn func(*game.Game)) bool {
	now := s.now()
	s.mu.Lock()
	s.prune(now)
	entry, ok := s.games[id]
	if ok {
		entry.lastAccess = now
	}
	s.mu.Unlock()
	if !ok {
		return false
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	fn(entry.game)
	return true
}
```

`Create` calls `prune`, stores a `session`, and returns the game ID. `prune` deletes entries whose `now.Sub(lastAccess) > sessionTTL`.

- [ ] **Step 4: Verify store tests GREEN**

Run: `go test ./api/store`

Expected: both tests pass.

- [ ] **Step 5: Add failing handler regressions**

Create `api/handler/handler_test.go`. Use `httptest.NewRequest`, `request.SetPathValue("id", g.ID)`, and `httptest.NewRecorder` to prove:

```go
func TestDealHidesDealerHoleCard(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
	store.Create(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/"+g.ID+"/deal", nil)
	req.SetPathValue("id", g.ID)
	res := httptest.NewRecorder()
	NewDeal(store)(res, req)
	var body struct{ DealerCard string `json:"dealerCard"` }
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil { t.Fatal(err) }
	if body.DealerCard == "" { t.Fatal("dealer upcard missing") }
	store.WithGame(g.ID, func(g *game.Game) {
		if len(g.Dealer.Hand.Cards) != 2 { t.Fatalf("dealer cards = %d", len(g.Dealer.Hand.Cards)) }
	})
}

func TestAbandonCompletesActiveRound(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
	g.StartRound()
	store.Create(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/"+g.ID+"/abandon", nil)
	req.SetPathValue("id", g.ID)
	res := httptest.NewRecorder()
	NewAbandon(store)(res, req)
	if res.Code != http.StatusNoContent { t.Fatalf("status = %d", res.Code) }
}
```

- [ ] **Step 6: Run handler tests and verify RED**

Run: `go test ./api/handler`

Expected: compile failure because `NewAbandon` is missing and handlers still call `Get`.

- [ ] **Step 7: Move complete handler operations under `WithGame`**

In deal, move, and hint handlers, replace manual path splitting and `Get` with `r.PathValue("id")` and:

```go
var response struct {
	PlayerCards []string `json:"playerCards"`
	DealerCard  string   `json:"dealerCard"`
}
if !store.WithGame(r.PathValue("id"), func(g *game.Game) {
	if g.RoundState == game.RoundStateActive {
		writeError(w, http.StatusConflict, "ROUND_ALREADY_ACTIVE", "Round is already active")
		return
	}
	g.StartRound()
	response = struct {
		PlayerCards []string `json:"playerCards"`
		DealerCard  string   `json:"dealerCard"`
	}{
		PlayerCards: helpers.CardsToStrings(g.Player.ActiveHand.Cards),
		DealerCard:  g.Dealer.Hand.Cards[0].ToString(),
	}
}) {
	writeError(w, http.StatusNotFound, "GAME_NOT_FOUND", "Game not found")
}
```

Create `NewAbandon`, returning `409 NO_ACTIVE_ROUND` unless the round is active, otherwise calling `g.AbandonRound()` and returning `204`.

Replace suffix routing in `server.Start` with exact native patterns:

```go
mux.HandleFunc("POST /api/game", handler.NewGame(s.sessionStore))
mux.HandleFunc("POST /api/game/{id}/deal", handler.NewDeal(s.sessionStore))
mux.HandleFunc("POST /api/game/{id}/move", handler.NewMove(s.sessionStore))
mux.HandleFunc("GET /api/game/{id}/hint", handler.NewHint(s.sessionStore, s.advisor))
mux.HandleFunc("POST /api/game/{id}/abandon", handler.NewAbandon(s.sessionStore))
```

Remove redundant method checks from handlers because `ServeMux` now owns method dispatch.

- [ ] **Step 8: Verify API tests and race checks**

Run: `go test ./api/...`

Run: `go test -race ./api/store ./api/handler`

Expected: all pass without a race report.

- [ ] **Step 9: Commit session and API fixes**

```bash
git add api/server.go api/store api/handler
git commit -m "fix: serialize game sessions"
```

### Task 3: Make the trainer follow server state safely

**Files:**
- Modify: `web/src/api.ts`
- Modify: `web/src/App.svelte`
- Create: `web/src/App.test.ts`
- Modify: `web/package.json`
- Modify: `web/package-lock.json`

**Interfaces:**
- Consumes: `POST /api/game/{id}/abandon` from Task 2.
- Produces: `abandonRound(): Promise<void>`, decision gating until a current hint exists, safe Next behavior, and rendering of all resolved hands.

- [ ] **Step 1: Add failing component regressions using Vitest and raw DOM APIs**

Mock `./api` before importing `App.svelte`:

```ts
const api = vi.hoisted(() => ({
  loadDeal: vi.fn(),
  makeMove: vi.fn(),
  getHint: vi.fn(),
  abandonRound: vi.fn()
}));
vi.mock('./api', () => api);

const flush = async () => {
  await Promise.resolve();
  await tick();
  await Promise.resolve();
  await tick();
};
```

Add three tests:

```ts
test('disables decisions until the current hint resolves', async () => {
  let resolveHint!: (value: string) => void;
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockReturnValue(new Promise(resolve => { resolveHint = resolve; }));
  const target = document.createElement('div');
  const app = new App({ target });
  await flush();
  const hit = [...target.querySelectorAll('button')].find(button => button.textContent?.startsWith('HIT'))!;
  expect(hit.disabled).toBe(true);
  resolveHint('HIT');
  await flush();
  expect(hit.disabled).toBe(false);
  app.$destroy();
});

test('abandons an active round before dealing after a wrong answer', async () => {
  api.loadDeal.mockResolvedValue({ playerCards: ['10', '6'], dealerCard: '10' });
  api.getHint.mockResolvedValue('HIT');
  api.makeMove.mockResolvedValue({ roundState: 'active', activeHand: ['10', '6'], unresolvedHands: [], resolvedHands: [], dealerCards: ['10'], outcomes: [] });
  api.abandonRound.mockResolvedValue(undefined);
  const target = document.createElement('div');
  const app = new App({ target });
  await flush();
  const stand = [...target.querySelectorAll('button')].find(button => button.textContent?.startsWith('STAND'))!;
  stand.click();
  await flush();
  const next = target.querySelector<HTMLButtonElement>('.next-btn')!;
  next.click();
  await flush();
  expect(api.abandonRound).toHaveBeenCalledOnce();
  expect(api.loadDeal).toHaveBeenCalledTimes(2);
  app.$destroy();
});

test('renders every resolved split hand', async () => {
  api.loadDeal.mockResolvedValue({ playerCards: ['8', '8'], dealerCard: '6' });
  api.getHint.mockResolvedValue('SPLIT');
  api.makeMove.mockResolvedValue({ roundState: 'complete', activeHand: [], unresolvedHands: [], resolvedHands: [['8', '10'], ['8', '9']], dealerCards: ['6', '10', '2'], outcomes: ['Push', 'Lose'] });
  const target = document.createElement('div');
  const app = new App({ target });
  await flush();
  const split = [...target.querySelectorAll('button')].find(button => button.textContent?.startsWith('SPLIT'))!;
  split.click();
  await flush();
  expect(target.querySelectorAll('.resolved-hand')).toHaveLength(2);
  app.$destroy();
});

```

Reset mocks and `document.body` in `beforeEach`.

- [ ] **Step 2: Run component tests and verify RED**

Run: `npm test --prefix web -- src/App.test.ts`

Expected: compile/test failures because `abandonRound`, hint gating, and `.resolved-hand` rendering do not exist.

- [ ] **Step 3: Implement API abandonment and hint-safe UI state**

Add to `api.ts`:

```ts
export async function abandonRound(): Promise<void> {
  if (!gameId) return;
  const res = await fetch(`/api/game/${gameId}/abandon`, { method: 'POST' });
  if (!res.ok) throw new Error('Failed to abandon round');
}
```

In `App.svelte`, add `hintLoading`, pass the hand key into `updateHint`, and only accept its result when that key is still current:

```ts
async function updateHint(key: string) {
  hintLoading = true;
  const nextHint = await getHint();
  if (hintKey === key) {
    hint = nextHint;
    hintLoading = false;
  }
}

$: canDecide = roundState === 'active' && !locked && !busy && !hintLoading && !!hint;
```

Use `canDecide` for all four buttons, keyboard decisions, and the guard in `decide`. When Next is pressed during an active round, call `await abandonRound()` before `handleLoadDeal()`.

Replace the single-hand completion branch with an `{#each resolvedHands as hand}` block carrying class `resolved-hand`. Remove the unused `outcomes` variable and assignments.

- [ ] **Step 4: Verify UI regressions GREEN**

Run: `npm test --prefix web -- src/App.test.ts`

Expected: all three tests pass. The first test also guarantees a hand cannot advance while its hint request is outstanding, preventing an older response from becoming current.

- [ ] **Step 5: Remove unused UI testing packages**

Run: `npm uninstall --prefix web --save-dev @testing-library/jest-dom @testing-library/svelte`

These tests use Vitest, Svelte, and native DOM APIs already present in the project.

- [ ] **Step 6: Run all UI tests and type checking**

Run: `npm test --prefix web`

Run: `npx tsc --noEmit -p web/tsconfig.json`

Expected: all pass.

- [ ] **Step 7: Commit UI fixes**

```bash
git add web/src web/package.json web/package-lock.json
git commit -m "fix: keep trainer rounds in sync"
```

### Task 4: Delete dead code and restore quality gates

**Files:**
- Modify: `api/helpers/convert.go`
- Modify: `internal/player/player.go`
- Modify: `internal/player/player_test.go`
- Modify: `internal/deck/deck.go`
- Modify: `internal/deck/deck_test.go`
- Modify: `internal/game/game.go`
- Modify: `web/eslint.config.js`
- Modify: `.golangci.yml`
- Modify: `README.md`

**Interfaces:**
- Consumes: the passing behavior from Tasks 1–3.
- Produces: a smaller public surface and working lint/build commands.

- [ ] **Step 1: Delete confirmed unused code**

Keep only `CardsToStrings` in `api/helpers/convert.go`. Remove the unused player-package error variables and `Player.CanMove`, plus their tests. Change `deck.NewDeck(rng)` to `deck.NewDeck()` because construction does not use randomness. In `generateShuffledShoe`, build unshuffled decks and shuffle only the combined shoe:

```go
for range decksInShoe {
	decks = append(decks, deck.NewDeck())
}
```

Update deck tests to call `NewDeck()`.

- [ ] **Step 2: Restore ESLint**

Add `alert: 'readonly'` to Svelte browser globals. The unused `outcomes` variable was removed in Task 3.

Run: `npm run lint --prefix web`

Expected: zero errors.

- [ ] **Step 3: Migrate golangci-lint configuration to v2**

Replace `.golangci.yml` with the validated v2 structure produced by the installed migrator:

```yaml
version: "2"
run:
  tests: true
linters:
  default: none
  enable: [errcheck, gocritic, govet, ineffassign, staticcheck, unparam, unused]
  settings:
    errcheck:
      check-type-assertions: true
      check-blank: false
    gocritic:
      enabled-tags: [diagnostic, experimental, opinionated, style]
      disabled-checks: [dupImport, ifElseChain, octalLiteral, whyNoLint, paramTypeCombine, unnamedResult]
    staticcheck:
      checks: [all]
  exclusions:
    rules:
      - path: _test\.go
        linters: [errcheck, gocritic]
      - path: api/handler/.*\.go
        text: Error return value of.*is not checked
        linters: [errcheck]
issues:
  max-issues-per-linter: 0
  max-same-issues: 0
formatters:
  enable: [gofmt, goimports]
```

Run: `golangci-lint config verify`

Run: `golangci-lint run ./...`

Expected: configuration verification and lint both pass.

- [ ] **Step 4: Correct README claims**

Replace “Winnings tracking with pot, bet, and round outcomes” with “Round outcome tracking across normal and split hands.” Add the supported-rule summary: six decks, S17, double after split, no surrender, no insurance.

- [ ] **Step 5: Run complete verification**

Run: `gofmt -w api/*.go api/handler/*.go api/store/*.go internal/*/*.go internal/strategy/strategies/*.go`

Run: `go test ./...`

Run: `go test -race ./api/store ./api/handler ./internal/game`

Run: `go vet ./...`

Run: `npm test --prefix web`

Run: `web/node_modules/.bin/tsc --noEmit -p web/tsconfig.json`

Run: `make lint`

Run: `make build`

Run: `git diff --check`

Expected: every command exits zero with no test, race, vet, type, lint, build, or whitespace failures.

- [ ] **Step 6: Commit cleanup**

```bash
git add .golangci.yml README.md api/helpers internal web/eslint.config.js
git commit -m "chore: remove dead code and repair checks"
```
