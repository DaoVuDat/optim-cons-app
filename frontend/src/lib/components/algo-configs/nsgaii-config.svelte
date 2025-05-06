<script lang="ts">
  import {nsgaiiConfig} from "$lib/stores/algorithms";

  let config = nsgaiiConfig
</script>

<div class="p-2 w-full h-full">
  <h2 class="text-2xl text-center font-bold">
    NSGA-II Algorithm
  </h2>
  <div class="grid gap-2 grid-cols-2 grid-rows-4">
    <fieldset class="fieldset text-lg text-nowrap">
      <legend class="fieldset-legend">Generations:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.generation} />
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg  text-nowrap">Chromosomes:</legend>
      <input type="number" class="input input-lg" placeholder="100" bind:value={config.chromosome}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg  text-nowrap">Crossover Rate:</legend>
      <input type="number" class="input input-lg" placeholder="0.9" bind:value={config.crossoverRate}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Mutation Rate:</legend>
      <input type="number" class="input input-lg" placeholder="0.1" bind:value={config.mutationRate}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Mutation Strength:</legend>
      <input type="number" class="input input-lg" placeholder="0.1" bind:value={config.mutationStrength}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Sigma:</legend>
      <input type="number" class="input input-lg" placeholder="2" bind:value={config.sigma}/>
    </fieldset>
  </div>
</div>
