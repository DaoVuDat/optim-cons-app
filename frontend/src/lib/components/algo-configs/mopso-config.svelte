<script lang="ts">
  import {mopsoConfig} from "$lib/stores/algorithms";

  let config = mopsoConfig
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
      <legend class="fieldset-legend text-lg text-nowrap">Mutation rate:</legend>
      <input type="number" class="input input-lg" placeholder="0.5" bind:value={config.mutationRate}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Archive Size:</legend>
      <input type="number" class="input input-lg" placeholder="100" bind:value={config.archiveSize}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Number of Grids:</legend>
      <input type="number" class="input input-lg" placeholder="20" bind:value={config.numberOfGrids}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">C1:</legend>
      <input type="number" class="input input-lg" placeholder="2" bind:value={config.c1}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">C2:</legend>
      <input type="number" class="input input-lg" placeholder="2" bind:value={config.c2}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">W:</legend>
      <input type="number" class="input input-lg" placeholder="0.4" bind:value={config.w}/>
    </fieldset>
  </div>
</div>
