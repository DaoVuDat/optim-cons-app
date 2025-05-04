<script lang="ts">

  import {gaConfig} from "$lib/stores/algorithms";

  let config = gaConfig
</script>


<div class="p-2 w-full h-full">
  <h2 class="text-2xl text-center font-bold">
    Genetic Algorithm
  </h2>
  <div class="grid gap-2 grid-cols-2 grid-rows-3">
    <fieldset class="fieldset text-lg text-nowrap">
      <legend class="fieldset-legend">Generations:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.generation} />
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Chromosomes:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.chromosome}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Crossover Rate:</legend>
      <input type="number" class="input input-lg" placeholder="0.7" bind:value={config.crossoverRate}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Mutation Rate:</legend>
      <input type="number" class="input input-lg" placeholder="0.1" bind:value={config.mutationRate}/>
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg text-nowrap">Elitism Count:</legend>
      <input type="number" class="input input-lg" placeholder="5" bind:value={config.elitismCount}/>
    </fieldset>
  </div>
</div>