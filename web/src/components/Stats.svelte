<script lang="ts">
  export let value: number;
  export let label: string;
</script>

<div class="stat">
  <div class="stat-value">{value}</div>
  <div class="stat-label">{label}</div>
</div>

<style>
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
</style>

