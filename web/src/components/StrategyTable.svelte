<script lang="ts">
  import { TABLE_HEADERS, HARD_MATRIX, SOFT_MATRIX, PAIR_MATRIX } from '../constants';

  const sections = [
    { title: 'Hard Hands', rows: HARD_MATRIX },
    { title: 'Soft Hands', rows: SOFT_MATRIX },
    { title: 'Pairs', rows: PAIR_MATRIX }
  ];
  const decisionClasses: Record<string, string> = {
    S: 'cell-stand', H: 'cell-hit', D: 'cell-double', SP: 'cell-split'
  };
</script>

<details class="info-wrap">
  <summary aria-label="Show strategy tables">i</summary>
  <div class="info-panel">
    {#each sections as section}
      <div class="section">{section.title}</div>
      <table>
        <tr>
          {#each TABLE_HEADERS as header}
            <th>{header}</th>
          {/each}
        </tr>
        {#each section.rows as row}
          <tr>
            <td>{row[0]}</td>
            {#each row.slice(1) as decision}
              <td class={decisionClasses[decision] ?? ''}>{decision}</td>
            {/each}
          </tr>
        {/each}
      </table>
    {/each}
  </div>
</details>

<style>
  .info-wrap {
    position: fixed;
    top: 16px;
    right: 16px;
    width: 36px;
    height: 36px;
    z-index: 10;
  }

  .info-wrap summary {
    width: 36px;
    height: 36px;
    list-style: none;
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

  .info-wrap summary::-webkit-details-marker {
    display: none;
  }

  .info-panel {
    width: 320px;
    background: rgba(0, 0, 0, 0.92);
    border: 2px solid #555;
    border-radius: 12px;
    padding: 16px;
    font-size: 12px;
    line-height: 1.4;
    position: absolute;
    top: 0;
    right: 44px;
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
