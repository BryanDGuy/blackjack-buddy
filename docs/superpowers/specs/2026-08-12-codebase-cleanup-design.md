# Codebase Cleanup Design

## Goal

Fix the reviewed correctness and state-management defects, restore all existing quality checks, and remove clearly dead code without redesigning the application.

## Blackjack Rules

- Six-deck shoe.
- Dealer stands on soft 17.
- Double down is allowed only on a hand's first two cards.
- Double after split is allowed.
- Surrender and insurance are not supported.
- A two-card 21 created by a split is a normal 21, not a natural blackjack.

## Game Engine

Deal two cards to both player and dealer at round start, while exposing only the dealer's upcard until the round completes. Track whether a hand came from a split so outcome calculation can distinguish split 21 from natural blackjack. Reject double down after a hit. Preserve the current decision matrix and six-deck shoe behavior.

The engine remains the authority for legal moves and outcomes. The UI may disable unavailable actions for usability, but the API must still reject illegal requests.

## Session and API State

Serialize all reads and mutations for an individual game. The session store will expose a callback-based operation that holds a per-game lock for the complete handler operation instead of returning an unlocked mutable pointer.

Expire idle sessions after 24 hours during store access and creation. Cleanup stays opportunistic; a background worker is unnecessary for this application.

Starting the next round after an incorrect answer must first finish the active round. Add an explicit abandon operation that completes the current round without scoring further decisions, then start a fresh deal. Keep endpoint behavior and response shapes otherwise stable.

Use exact Go `ServeMux` method/path patterns and path values instead of suffix routing and repeated manual path splitting.

## Web Client

Treat hint loading as part of the busy state. Decisions remain disabled until the hint for the current hand has arrived, and stale responses cannot replace a newer hand's hint. After an incorrect non-terminal decision, Next abandons the server round before dealing again.

Render every resolved split hand at round completion. Remove unused local state and unused testing-library packages. Keep the existing visual design and keyboard controls.

## Cleanup and Tooling

Remove unused card parsing/conversion functions and redundant per-deck shuffling. Update the golangci-lint configuration to the installed v2 format. Add `alert` to browser globals and remove the unused UI variable so the current lint policy passes.

Do not add dependencies, abstractions, persistence, betting, surrender, insurance, or a broader UI redesign.

## Testing

Add focused regression tests before implementation for:

- double down after a hit;
- split 21 outcome classification;
- dealer initial card count and hidden-card API response;
- serialized access to one game;
- idle-session expiry;
- abandoning an active round;
- hint gating/stale response handling and completed split-hand rendering.

Verification requires Go tests, race-enabled Go tests where supported, Go vet, UI tests, TypeScript checking, UI lint, Go lint, and a production build.
