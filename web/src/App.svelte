<script lang="ts">
  import { onMount } from 'svelte';
  import { handleLoadDeal, decide, startNextRound, handleKey, updateHint, type GameState } from './game';
  import Stats from './components/Stats.svelte';
  import DeckPanel from './components/DeckPanel.svelte';
  import StrategyTable from './components/StrategyTable.svelte';

  const state: GameState = {
    skipTrivial: false,
    playerCards: [],
    dealerCard: '',
    dealerCards: [],
    queuedHands: [],
    completedHands: [],
    correct: 0,
    total: 0,
    pot: 1000,
    bet: 10,
    totalWinnings: 0,
    roundWinnings: 0,
    deckState: { totalCards: 0, rankCounts: {} },
    hint: '',
    resultText: 'Result: -',
    resultClass: 'result',
    outcomeText: 'Outcome: -',
    outcomeClass: 'result outcome-box',
    nextVisible: false,
    busy: false,
    locked: false
  };

  $: percent = state.total > 0 ? Math.round((state.correct / state.total) * 100) : 0;
  $: derivedDealerCards = state.dealerCards.length ? state.dealerCards : state.dealerCard ? [state.dealerCard] : [];
  
  let hintKey = '';
  $: {
    const newKey = `${state.playerCards.join(',')}-${state.dealerCard}`;
    if (newKey !== hintKey && state.playerCards.length && state.dealerCard && !state.locked) {
      hintKey = newKey;
      updateHint(state);
    }
  }

  onMount(() => {
    handleLoadDeal(state);
    const keyHandler = (e: KeyboardEvent) => handleKey(state, e);
    window.addEventListener('keydown', keyHandler);
    return () => window.removeEventListener('keydown', keyHandler);
  });
</script>

<div class="app">
  <div class="stats">
    <Stats value={state.correct} label="CORRECT" />
    <Stats value={state.total} label="TOTAL" />
    <Stats value={percent} label="ACCURACY" />
  </div>

  <div class="stats">
    <Stats value={state.pot} label="POT" />
    <Stats value={state.bet} label="BET" />
    <div class="stat">
      <div class="stat-value">{state.totalWinnings >= 0 ? '+' : ''}{state.totalWinnings}</div>
      <div class="stat-label">WINNINGS</div>
    </div>
  </div>

  <div class="options">
    <label class="checkbox"><input type="checkbox" bind:checked={state.skipTrivial} on:change={() => handleLoadDeal(state)}> Skip Trivial</label>
  </div>

  <div class="section">
    <div class="label">Dealer Card</div>
    <div class="cards">
      {#each derivedDealerCards as card, index}
        <div class="card" data-index={index}>{card}</div>
      {/each}
    </div>
  </div>

  <div class="section">
    <div class="label">Your Cards</div>
    <div class="cards">
      {#each state.playerCards as card, index}
        <div class="card" data-index={index}>{card}</div>
      {/each}
    </div>
  </div>

  <div class="buttons-section">
    <div class="buttons">
      <button on:click={() => decide(state, 'HIT')} disabled={state.locked || state.busy}>HIT (H)</button>
      <button on:click={() => decide(state, 'STAND')} disabled={state.locked || state.busy}>STAND (S)</button>
      <button on:click={() => decide(state, 'DOUBLE DOWN')} disabled={state.locked || state.busy}>DOUBLE DOWN (D)</button>
      <button on:click={() => decide(state, 'SPLIT')} disabled={state.locked || state.busy}>SPLIT (P)</button>
    </div>
    {#if state.hint}
      <div class="hint-wrap">
        <div class="hint-button">?</div>
        <div class="hint-panel hint-{state.hint.toLowerCase().replace(/\s+/g, '-')}">Hint: {state.hint}</div>
      </div>
    {/if}
  </div>

  <div class={state.resultClass}>{state.resultText}</div>
  <div class={state.outcomeClass}>{state.outcomeText}</div>
  <button class="next-btn" class:visible={state.nextVisible} on:click={() => startNextRound(state)}>Next (N)</button>

  <DeckPanel deckState={state.deckState} />
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

  .stat {
    text-align: center;
  }

  .stat-value {
    font-size: 28px;
    font-weight: bold;
  }

  .stat-label {
    font-size: 12px;
    color: #888;
    margin-top: 4px;
  }

  .options {
    text-align: center;
  }

  .checkbox {
    font-size: 14px;
    color: #fff;
    cursor: pointer;
  }

  .checkbox input {
    margin-right: 8px;
    cursor: pointer;
  }

  .section {
    margin: 16px 0;
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

  button:hover {
    background: #444;
    border-color: #777;
  }

  button:active {
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
    min-height: 64px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .next-btn {
    margin-top: 20px;
    width: 100%;
    padding: 14px;
    background: #2d2d2d;
    color: #fff;
    border: 2px solid #555;
    border-radius: 8px;
    font-size: 16px;
    cursor: pointer;
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
