<script>
  import { onMount } from 'svelte';
  import {
    KEY_BINDINGS,
    TABLE_HEADERS,
    HARD_MATRIX,
    SOFT_MATRIX,
    PAIR_MATRIX
  } from './constants.js';

  let skipTrivial = true;
  let playerCards = [];
  let dealerCard = '';
  let dealerCards = [];
  let queuedHands = [];
  let completedHands = [];
  let completedOutcomes = [];
  let correct = 0;
  let total = 0;
  let resultText = 'Result: -';
  let resultClass = 'result';
  let outcomeText = 'Outcome: -';
  let outcomeClass = 'result outcome-box';
  let nextVisible = false;

  $: percent = total > 0 ? Math.round((correct / total) * 100) : 0;
  $: derivedDealerCards = dealerCards.length ? dealerCards : dealerCard ? [dealerCard] : [];

  const rowDecisions = (row) => row.slice(1);

  const decisionClass = (val) => {
    switch (val) {
      case 'S':
        return 'cell-stand';
      case 'H':
        return 'cell-hit';
      case 'D':
        return 'cell-double';
      case 'SP':
        return 'cell-split';
      default:
        return '';
    }
  };

  async function loadScenario() {
    const res = await fetch('/api/scenario', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ skipTrivial })
    });
    const data = await res.json();
    playerCards = data.playerCards || [];
    dealerCard = data.dealerCard || '';
    dealerCards = dealerCard ? [dealerCard] : [];
    queuedHands = [];
    completedHands = [];
    completedOutcomes = [];
    resultText = 'Result: -';
    resultClass = 'result';
    outcomeText = 'Outcome: -';
    outcomeClass = 'result outcome-box';
    nextVisible = false;
  }

  async function decide(decision) {
    const res = await fetch('/api/check', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        playerCards,
        dealerCard,
        decision,
        queuedHands,
        completedHands
      })
    });
    const result = await res.json();

    total += 1;
    if (result.correct) {
      correct += 1;
    }

    resultClass = `result ${result.correct ? 'correct' : 'incorrect'}`;
    resultText = `${result.correct ? 'Correct' : 'Incorrect'}: ${result.correctDecision}`;
    outcomeClass = 'result outcome-box';
    outcomeText = result.outcome ? `Outcome: ${result.outcome}` : 'Outcome: -';

    if (!result.correct || result.restart) {
      nextVisible = true;
    }

    if (Array.isArray(result.playerCards) && result.playerCards.length) {
      playerCards = result.playerCards;
    }
    if (Array.isArray(result.dealerCards) && result.dealerCards.length) {
      dealerCards = result.dealerCards;
    }

    queuedHands = Array.isArray(result.queuedHands) ? result.queuedHands : [];
    completedHands = Array.isArray(result.completedHands) ? result.completedHands : [];
    completedOutcomes = Array.isArray(result.completedOutcomes) ? result.completedOutcomes : [];

    if (result.roundComplete && queuedHands.length === 0 && result.correct && !result.restart) {
      nextVisible = true;
    }
    if (!result.roundComplete && result.correct && !result.restart) {
      nextVisible = false;
    }
  }

  function handleKey(event) {
    const tag = event.target?.tagName;
    if (tag && ['INPUT', 'SELECT', 'TEXTAREA'].includes(tag)) {
      return;
    }
    const action = KEY_BINDINGS[event.key.toLowerCase()];
    if (!action) return;
    event.preventDefault();
    if (action === 'NEXT') {
      if (nextVisible) {
        loadScenario();
      }
      return;
    }
    decide(action);
  }

  onMount(() => {
    loadScenario();
    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  });
</script>

<div class="app">
  <div class="stats">
    <div class="stat"><div class="stat-value">{correct}</div><div class="stat-label">CORRECT</div></div>
    <div class="stat"><div class="stat-value">{total}</div><div class="stat-label">TOTAL</div></div>
    <div class="stat"><div class="stat-value">{percent}%</div><div class="stat-label">ACCURACY</div></div>
  </div>

  <div class="options">
    <label class="checkbox"><input type="checkbox" bind:checked={skipTrivial} on:change={loadScenario}> Skip Trivial</label>
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
      {#each playerCards as card, index}
        <div class="card" data-index={index}>{card}</div>
      {/each}
    </div>
  </div>

  <div class="buttons">
    <button on:click={() => decide('HIT')}>HIT (H)</button>
    <button on:click={() => decide('STAND')}>STAND (S)</button>
    <button on:click={() => decide('DOUBLE DOWN')}>DOUBLE DOWN (D)</button>
    <button on:click={() => decide('SPLIT')}>SPLIT (P)</button>
  </div>

  <div class={resultClass}>{resultText}</div>
  <div class={outcomeClass}>{outcomeText}</div>
  <button class="next-btn" class:visible={nextVisible} on:click={loadScenario}>Next (N)</button>

  <div class="info-wrap">
    <div class="info-button">i</div>
    <div class="info-panel">
      <div class="section">Hard Hands</div>
      <table>
        <tr>
          {#each TABLE_HEADERS as header}
            <th>{header}</th>
          {/each}
        </tr>
        {#each HARD_MATRIX as row}
          <tr>
            <td>{row[0]}</td>
            {#each rowDecisions(row) as decision}
              <td class={decisionClass(decision)}>{decision}</td>
            {/each}
          </tr>
        {/each}
      </table>
      <div class="section">Soft Hands</div>
      <table>
        <tr>
          {#each TABLE_HEADERS as header}
            <th>{header}</th>
          {/each}
        </tr>
        {#each SOFT_MATRIX as row}
          <tr>
            <td>{row[0]}</td>
            {#each rowDecisions(row) as decision}
              <td class={decisionClass(decision)}>{decision}</td>
            {/each}
          </tr>
        {/each}
      </table>
      <div class="section">Pairs</div>
      <table>
        <tr>
          {#each TABLE_HEADERS as header}
            <th>{header}</th>
          {/each}
        </tr>
        {#each PAIR_MATRIX as row}
          <tr>
            <td>{row[0]}</td>
            {#each rowDecisions(row) as decision}
              <td class={decisionClass(decision)}>{decision}</td>
            {/each}
          </tr>
        {/each}
      </table>
    </div>
  </div>
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

  .buttons {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
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

  .info-wrap {
    position: fixed;
    top: 16px;
    right: 16px;
    display: flex;
    gap: 12px;
    align-items: flex-start;
    flex-direction: row-reverse;
    z-index: 10;
  }

  .info-button {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: 2px solid #555;
    background: rgba(0, 0, 0, 0.4);
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    font-weight: bold;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
  }

  .info-panel {
    width: 320px;
    background: rgba(0, 0, 0, 0.92);
    border: 2px solid #555;
    border-radius: 12px;
    padding: 16px;
    font-size: 12px;
    line-height: 1.4;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
  }

  .info-panel table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 12px;
  }

  .info-panel th,
  .info-panel td {
    padding: 4px;
    text-align: center;
    border: 1px solid rgba(255, 255, 255, 0.1);
    font-size: 11px;
  }

  .info-panel .section {
    font-weight: bold;
    text-align: left;
    padding: 6px 0 4px;
  }

  .info-wrap:hover .info-panel {
    opacity: 1;
    pointer-events: auto;
  }

  .cell-stand {
    background: #7b1f1f;
    color: #fff;
  }

  .cell-hit {
    background: #2d4f2d;
    color: #fff;
  }

  .cell-double {
    background: #7b6b1f;
    color: #fff;
  }

  .cell-split {
    background: #1f3f7b;
    color: #fff;
  }
</style>
