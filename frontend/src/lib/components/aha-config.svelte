<script lang="ts">
  import {Algorithms, algorithmsStore} from "$lib/stores/algorithms.svelte";

  let config = algorithmsStore.getConfig(Algorithms.AHA)
</script>


<div class="p-2 w-full h-full">
  <h2 class="text-2xl text-center font-bold">
    Artificial Hummingbird Algorithm
  </h2>
  <div class="grid gap-2 grid-cols-1 grid-rows-3">
    <fieldset class="fieldset text-lg">
      <legend class="fieldset-legend">Iterations:</legend>
      <input type="text" class="input input-lg" placeholder="300" bind:value={config.iterations.value} />
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg">Populations:</legend>
      <input type="text" class="input input-lg" placeholder="300" bind:value={config.population.value}/>
    </fieldset>
  </div>
</div>