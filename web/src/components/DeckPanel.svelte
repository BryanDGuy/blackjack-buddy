<script lang="ts">
  export let shoeState: { totalCards: number; rankCounts: Record<string, number> };
</script>

<div class="deck-wrap">
  <div class="deck-button">D</div>
  <div class="deck-panel">
    <div class="section">Deck: {shoeState.totalCards} cards</div>
    <div class="deck-ranks">
      {#each (() => {
        const entries = Object.entries(shoeState.rankCounts || {});
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

<style>
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
</style>

