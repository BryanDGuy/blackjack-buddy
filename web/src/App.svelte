<script>
  import { onMount } from 'svelte';
  import {
    KEY_BINDINGS,
    TABLE_HEADERS,
    HARD_MATRIX,
    SOFT_MATRIX,
    PAIR_MATRIX
  } from './constants.js';

  let skipTrivial = false;
  let playerCards = [];
  let dealerCard = '';
  let dealerCards = [];
  let queuedHands = [];
  let completedHands = [];
  let completedOutcomes = [];
  let correct = 0;
  let total = 0;
  let pot = 1000;
  let bet = 10;
  let totalWinnings = 0;
  let roundWinnings = 0;
  let deckState = { totalCards: 0, rankCounts: {} };
  let hint = '';
  let resultText = 'Result: -';
  let resultClass = 'result';
  let outcomeText = 'Outcome: -';
  let outcomeClass = 'result outcome-box';
  let nextVisible = false;
  let busy = false;
  let locked = false;

  $: percent = total > 0 ? Math.round((correct / total) * 100) : 0;
  $: derivedDealerCards = dealerCards.length ? dealerCards : dealerCard ? [dealerCard] : [];
  
  let hintKey = '';
  $: {
    const newKey = `${playerCards.join(',')}-${dealerCard}`;
    if (newKey !== hintKey && playerCards.length && dealerCard && !locked) {
      hintKey = newKey;
      updateHint();
    }
  }
  
  async function updateHint() {
    if (!playerCards.length || !dealerCard || locked) return;
    try {
      const res = await fetch('/api/hint', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ playerCards, dealerCard })
      });
      const data = await res.json();
      hint = data.hint || '';
    } catch (e) {
      hint = '';
    }
  }

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
    if (total === 0) pot = data.pot || 1000;
    bet = data.bet || 10;
    roundWinnings = 0;
    if (data.deckState) deckState = data.deckState;
    hint = data.hint || '';
    resultText = 'Result: -';
    resultClass = 'result';
    outcomeText = 'Outcome: -';
    outcomeClass = 'result outcome-box';
    nextVisible = false;
    locked = false;
  }

  async function decide(decision) {
    if (locked || busy) {
      return;
    }

    busy = true;

    try {
      const res = await fetch('/api/check', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          playerCards,
          dealerCard,
          decision,
          queuedHands,
          completedHands,
          pot,
          bet,
          totalWinnings
        })
      });
      const result = await res.json();

      total += 1;
      if (result.correct) {
        correct += 1;
      }

      if (result.pot !== undefined) pot = result.pot;
      if (result.bet !== undefined) bet = result.bet;
      if (result.roundWinnings !== undefined) roundWinnings = result.roundWinnings;
      if (result.totalWinnings !== undefined) totalWinnings = result.totalWinnings;
      if (result.deckState) deckState = result.deckState;
      if (result.hint !== undefined) hint = result.hint;

      resultClass = `result ${result.correct ? 'correct' : 'incorrect'}`;
      resultText = `${result.correct ? 'Correct' : 'Incorrect'}: ${result.correctDecision}`;
      outcomeClass = 'result outcome-box';
      if (result.outcome) {
        let winningsText = '';
        if (result.roundComplete && result.roundWinnings !== undefined) {
          const sign = result.roundWinnings >= 0 ? '+' : '';
          winningsText = ` (${sign}${result.roundWinnings})`;
        }
        outcomeText = `Outcome: ${result.outcome}${winningsText}`;
      } else if (result.roundComplete && result.roundWinnings !== undefined) {
        const sign = result.roundWinnings >= 0 ? '+' : '';
        outcomeText = `Outcome: (${sign}${result.roundWinnings})`;
      } else {
        outcomeText = 'Outcome: -';
      }

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
    } finally {
      locked = nextVisible;
      busy = false;
    }
  }

  async function startNextRound() {
    if (!nextVisible || busy) {
      return;
    }
    busy = true;
    try {
      await loadScenario();
    } finally {
      busy = false;
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
      startNextRound();
      return;
    }
    if (locked || busy) {
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

  <div class="stats">
    <div class="stat"><div class="stat-value">{pot}</div><div class="stat-label">POT</div></div>
    <div class="stat"><div class="stat-value">{bet}</div><div class="stat-label">BET</div></div>
    <div class="stat"><div class="stat-value">{totalWinnings >= 0 ? '+' : ''}{totalWinnings}</div><div class="stat-label">WINNINGS</div></div>
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

  <div class="buttons-section">
    <div class="buttons">
      <button on:click={() => decide('HIT')} disabled={locked || busy}>HIT (H)</button>
      <button on:click={() => decide('STAND')} disabled={locked || busy}>STAND (S)</button>
      <button on:click={() => decide('DOUBLE DOWN')} disabled={locked || busy}>DOUBLE DOWN (D)</button>
      <button on:click={() => decide('SPLIT')} disabled={locked || busy}>SPLIT (P)</button>
    </div>
    {#if hint}
      <div class="hint-wrap">
        <div class="hint-button">?</div>
        <div class="hint-panel hint-{hint.toLowerCase().replace(/\s+/g, '-')}">Hint: {hint}</div>
      </div>
    {/if}
  </div>

  <div class={resultClass}>{resultText}</div>
  <div class={outcomeClass}>{outcomeText}</div>
  <button class="next-btn" class:visible={nextVisible} on:click={startNextRound}>Next (N)</button>

  <div class="deck-wrap">
    <div class="deck-button">D</div>
    <div class="deck-panel">
      <div class="section">Deck: {deckState.totalCards} cards</div>
      <div class="deck-ranks">
        {#each (() => {
          const entries = Object.entries(deckState.rankCounts || {});
          const order = ['A', '2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K'];
          return entries.sort((a, b) => order.indexOf(a[0]) - order.indexOf(b[0]));
        })() as [rank, count]}
          <div class="deck-rank-item">
            <span class="deck-rank">{rank}</span>
            <span class="deck-count">{count}</span>
          </div>
        {/each}
      </div>
    </div>
  </div>

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

  .deck-wrap {
    position: fixed;
    top: 16px;
    left: 16px;
    width: 36px;
    height: 36px;
    z-index: 10;
  }

  .deck-button {
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
    position: relative;
    z-index: 2;
  }

  .deck-panel {
    width: 200px;
    background: rgba(0, 0, 0, 0.92);
    border: 2px solid #555;
    border-radius: 12px;
    padding: 16px;
    font-size: 12px;
    line-height: 1.4;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
    position: absolute;
    top: 0;
    left: 44px;
    visibility: hidden;
  }

  .deck-wrap:hover .deck-panel {
    opacity: 1;
    pointer-events: auto;
    visibility: visible;
  }

  .deck-panel .section {
    font-weight: bold;
    text-align: left;
    padding: 6px 0 4px;
    margin-bottom: 8px;
  }

  .deck-ranks {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 8px;
    margin-top: 8px;
  }

  .deck-rank-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 4px;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
  }

  .deck-rank {
    font-weight: bold;
    font-size: 14px;
    margin-bottom: 2px;
  }

  .deck-count {
    font-size: 11px;
    color: #888;
  }

  .info-wrap {
    position: fixed;
    top: 16px;
    right: 16px;
    width: 36px;
    height: 36px;
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
    position: relative;
    z-index: 2;
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
    position: absolute;
    top: 0;
    right: 44px;
    visibility: hidden;
  }

  .info-wrap:hover .info-panel {
    opacity: 1;
    pointer-events: auto;
    visibility: visible;
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
