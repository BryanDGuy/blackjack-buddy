<script lang="ts">
  import { TABLE_HEADERS, HARD_MATRIX, SOFT_MATRIX, PAIR_MATRIX } from '../constants';
  import { rowDecisions, decisionClass } from '../utils';
</script>

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

<style>
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

