<script lang="ts">
    import {outOfBoundConfig} from "$lib/stores/constraints";

    const config = outOfBoundConfig
</script>

<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-1 grid-rows-2 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Power Difference (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="1" bind:value={config.PowerDifferencePenalty}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="20000" bind:value={config.AlphaOutOfBoundaryPenalty}/>
    </fieldset>
  </div>
</div>