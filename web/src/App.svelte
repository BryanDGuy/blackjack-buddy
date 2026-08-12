<script lang="ts">
  import { onMount } from 'svelte';
  import { KEY_BINDINGS } from './constants';
  import { abandonRound, loadDeal, makeMove, getHint } from './api';
  import Stats from './components/Stats.svelte';
  import StrategyTable from './components/StrategyTable.svelte';

  let playerCards: string[] = [];
  let dealerCards: string[] = [];
  let unresolvedHands: string[][] = [];
  let resolvedHands: string[][] = [];
  let correct = 0;
  let total = 0;
  let hint = '';
  let resultText = 'Result: -';
  let resultClass = 'result';
  let outcomeText = 'Outcome: -';
  let roundState = 'none';
  let nextVisible = false;
  let busy = false;
  let hintLoading = false;

  $: percent = total > 0 ? Math.round((correct / total) * 100) : 0;
  $: canDecide = roundState === 'active' && !nextVisible && !busy && !hintLoading && Boolean(hint);
  $: canDouble = canDecide && playerCards.length === 2;
  $: canSplit = canDouble && playerCards[0] === playerCards[1];

  let hintKey = '';
  let hintGeneration = 0;
  $: {
    if (!busy && roundState === 'active' && !nextVisible) {
      const newKey = `${hintGeneration}-${playerCards.join(',')}-${dealerCards[0]}`;
      if (newKey !== hintKey && playerCards.length && dealerCards[0]) {
        hintKey = newKey;
        updateHint(newKey);
      }
    } else if (roundState !== 'active') {
      hint = '';
      hintKey = '';
    }
  }
  
  async function updateHint(key: string) {
    if (!playerCards.length || !dealerCards[0] || nextVisible || roundState !== 'active' || busy) {
      return;
    }
    hintLoading = true;
    const nextHint = await getHint();
    if (hintKey === key) {
      hintLoading = false;
      if (nextHint) {
        hint = nextHint;
      } else {
        nextVisible = true;
        resultText = 'Hint unavailable. Start next round.';
      }
    }
  }

  async function handleLoadDeal() {
    try {
      const data = await loadDeal();
      playerCards = data.playerCards;
      dealerCards = data.dealerCard ? [data.dealerCard] : [];
      unresolvedHands = [];
      resolvedHands = [];
      roundState = 'active';
      hint = '';
      hintKey = '';
      hintGeneration += 1;
      resultText = 'Result: -';
      resultClass = 'result';
      outcomeText = 'Outcome: -';
      nextVisible = false;
    } catch (err) {
      console.error('Failed to load deal:', err);
      alert(`Failed to load deal: ${err instanceof Error ? err.message : String(err)}`);
    }
  }

  async function decide(decision: string) {
    if (!canDecide || (decision === 'DOUBLE DOWN' && !canDouble) || (decision === 'SPLIT' && !canSplit)) {
      return;
    }

    const expectedHint = hint;
    const isCorrect = expectedHint === decision;

    busy = true;

    try {
      const result = await makeMove(decision);

      total += 1;
      roundState = result.roundState;
      playerCards = result.activeHand;
      unresolvedHands = result.unresolvedHands;
      resolvedHands = result.resolvedHands || [];
      hintGeneration += 1;

      if (result.dealerCards.length > 0) {
        dealerCards = result.dealerCards;
      }

      if (isCorrect) {
        correct += 1;
      }

      resultClass = `result ${isCorrect ? 'correct' : 'incorrect'}`;
      resultText = `${isCorrect ? 'Correct' : 'Incorrect'}: ${expectedHint || decision}`;
      
      outcomeText = result.outcomes.length ? `Outcome: ${result.outcomes.join(' | ')}` : 'Outcome: -';
      nextVisible = result.roundState === 'complete' || !isCorrect;
    } catch (err) {
      console.error('Failed to make move:', err);
    } finally {
      busy = false;
    }
  }

  async function startNextRound() {
    if (!nextVisible || busy) {
      return;
    }
    busy = true;
    try {
      if (roundState === 'active') {
        await abandonRound();
        roundState = 'none';
      }
      await handleLoadDeal();
    } catch (err) {
      console.error('Failed to start next round:', err);
      alert(`Failed to start next round: ${err instanceof Error ? err.message : String(err)}`);
    } finally {
      busy = false;
    }
  }

  function handleKey(event: KeyboardEvent) {
    const target = event.target as HTMLElement;
    const tag = target?.tagName;
    if (tag && ['INPUT', 'SELECT', 'TEXTAREA'].includes(tag)) {
      return;
    }
    const action = KEY_BINDINGS[event.key.toLowerCase()];
    if (!action) return;
    event.preventDefault();
    if (action === 'NEXT') {
      startNextRound();
      return;
    }
    decide(action);
  }

  onMount(() => {
    handleLoadDeal();
    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  });
</script>

<div class="app">
  <div class="stats">
    <Stats value={correct} label="CORRECT" />
    <Stats value={total} label="TOTAL" />
    <Stats value={percent} label="ACCURACY" />
  </div>

  <div class="section">
    <div class="label">Dealer Card</div>
    <div class="cards">
      {#each dealerCards as card, index}
        <div class="card" data-index={index}>{card}</div>
      {/each}
    </div>
  </div>

  <div class="hands-container">
    <div class="section">
      <div class="label">Your Cards</div>
      {#if playerCards.length > 0}
        <div class="cards">
          {#each playerCards as card, index}
            <div class="card" data-index={index}>{card}</div>
          {/each}
        </div>
      {:else if roundState === 'complete'}
        {#each resolvedHands as hand}
          <div class="cards resolved-hand">
            {#each hand as card, index}
              <div class="card" data-index={index}>{card}</div>
            {/each}
          </div>
        {/each}
      {/if}
    </div>
    {#if unresolvedHands.length > 0}
      <div class="inactive-hands">
        {#each unresolvedHands as hand, handIndex}
          <div class="inactive-hand">
            <div class="label">Hand {handIndex + 1}</div>
            <div class="cards">
              {#each hand as card, index}
                <div class="card inactive" data-index={index}>{card}</div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <div class="buttons-section">
    <div class="buttons">
      <button on:click={() => decide('HIT')} disabled={!canDecide}>HIT (H)</button>
      <button on:click={() => decide('STAND')} disabled={!canDecide}>STAND (S)</button>
      <button on:click={() => decide('DOUBLE DOWN')} disabled={!canDouble}>DOUBLE DOWN (D)</button>
      <button on:click={() => decide('SPLIT')} disabled={!canSplit}>SPLIT (P)</button>
    </div>
    {#if hint}
      <div class="hint-wrap">
        <div class="hint-button">?</div>
        <div class="hint-panel hint-{hint.toLowerCase().replace(/\s+/g, '-')}">Hint: {hint}</div>
      </div>
    {/if}
  </div>

  <div class={resultClass}>{resultText}</div>
  <div class="result outcome-box">{outcomeText}</div>
  <button class="next-btn" class:visible={nextVisible} on:click={startNextRound}>Next (N)</button>

  <StrategyTable />
</div>

<style>
  :global(body) {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    background: #1a1a1a;
    color: #fff;
  }

  .app {
    width: 100%;
    max-width: 600px;
    padding: 20px;
    background: rgba(0, 0, 0, 0.35);
    border-radius: 16px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.35);
    margin: 40px auto;
    position: relative;
    min-height: calc(100vh - 80px);
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .stats {
    display: flex;
    justify-content: space-between;
    padding: 16px;
    background: #222;
    border-radius: 8px;
  }

  .hands-container {
    display: flex;
    gap: 20px;
    align-items: flex-start;
  }

  .section {
    margin: 16px 0;
    flex: 1;
  }

  .label {
    font-size: 14px;
    color: #888;
    margin-bottom: 12px;
    text-transform: uppercase;
    letter-spacing: 1px;
    text-align: center;
  }

  .cards {
    text-align: center;
    margin: 16px 0;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-wrap: nowrap;
  }

  .card {
    width: 80px;
    height: 112px;
    background: #fff;
    border-radius: 8px;
    margin: 0 8px;
    font-size: 32px;
    color: #000;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
  }

  .card.inactive {
    opacity: 0.4;
    background: #ddd;
  }

  .inactive-hands {
    display: flex;
    flex-direction: column;
    gap: 12px;
    min-width: 200px;
  }

  .inactive-hand {
    opacity: 0.5;
  }

  .inactive-hand .label {
    font-size: 12px;
    margin-bottom: 8px;
  }

  .inactive-hand .cards {
    margin: 8px 0;
  }

  .inactive-hand .card {
    width: 60px;
    height: 84px;
    font-size: 24px;
    margin: 0 4px;
  }

  .buttons-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .buttons {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
  }

  .hint-wrap {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .hint-button {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    border: 2px solid #555;
    background: rgba(0, 0, 0, 0.4);
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-weight: bold;
    font-size: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
    position: relative;
    z-index: 2;
  }

  .hint-panel {
    position: absolute;
    top: 0;
    left: 50%;
    transform: translateX(-50%) translateY(-100%);
    margin-bottom: 8px;
    padding: 8px 12px;
    border: 2px solid;
    border-radius: 8px;
    font-size: 12px;
    color: #fff;
    white-space: nowrap;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
    visibility: hidden;
    z-index: 10;
  }

  .hint-button:hover + .hint-panel {
    opacity: 1;
    pointer-events: auto;
    visibility: visible;
  }

  .hint-hit {
    background: #2d4f2d;
    border-color: #2d4f2d;
  }

  .hint-stand {
    background: #7b1f1f;
    border-color: #7b1f1f;
  }

  .hint-double-down {
    background: #7b6b1f;
    border-color: #7b6b1f;
  }

  .hint-split {
    background: #1f3f7b;
    border-color: #1f3f7b;
  }

  button {
    padding: 16px;
    background: #333;
    border: 2px solid #555;
    border-radius: 8px;
    color: #fff;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.2s;
  }

  button:not(:disabled):hover {
    background: #444;
    border-color: #777;
  }

  button:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  button:not(:disabled):active {
    transform: scale(0.98);
  }

  .result {
    padding: 20px;
    border-radius: 8px;
    text-align: center;
    font-size: 18px;
    min-height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .correct {
    background: #1a472a;
    border: 2px solid #2d8659;
  }

  .incorrect {
    background: #4a1a1a;
    border: 2px solid #862d2d;
  }

  .outcome-box {
    background: #222;
    border: 2px solid #555;
  }

  .next-btn {
    margin-top: 20px;
    width: 100%;
    padding: 14px;
    background: #2d2d2d;
    visibility: hidden;
    pointer-events: none;
    transition: background 0.2s ease;
  }

  .next-btn.visible {
    visibility: visible;
    pointer-events: auto;
  }

  .next-btn:hover {
    background: #3a3a3a;
  }

</style>
