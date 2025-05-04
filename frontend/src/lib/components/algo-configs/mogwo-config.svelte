<script lang="ts">
  import {mogwoConfig} from "$lib/stores/algorithms";

  let config = mogwoConfig
</script>

<div class="p-2 w-full h-full">
  <h2 class="text-2xl text-center font-bold">
    Multi-Objective Grey Wolf Optimizer
  </h2>
  <div class="grid gap-2 grid-cols-2 grid-rows-4">
    <fieldset class="fieldset text-lg">
      <legend class="fieldset-legend text-nowrap">Iterations:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.iterations} />
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Population:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.population}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">a Parameter:</legend>
      <input type="number" class="input input-lg" placeholder="2" bind:value={config.aParam}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Archive Size:</legend>
      <input type="number" class="input input-lg" placeholder="100" bind:value={config.archiveSize}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Number of Grids:</legend>
      <input type="number" class="input input-lg" placeholder="10" bind:value={config.numberOfGrids}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Alpha (Grid Inflation):</legend>
      <input type="number" class="input input-lg" placeholder="0.1" bind:value={config.alpha}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Beta (Leader Selection):</legend>
      <input type="number" class="input input-lg" placeholder="4" bind:value={config.beta}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Gamma (Repository Selection):</legend>
      <input type="number" class="input input-lg" placeholder="2" bind:value={config.gamma}/>
    </fieldset>
  </div>
</div>
