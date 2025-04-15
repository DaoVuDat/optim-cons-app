<script lang="ts">
  import {moahaConfig} from "$lib/stores/algorithms";

  let config = moahaConfig

</script>

<div class="p-2 w-full h-full">
  <h2 class="text-2xl text-center font-bold">
    Multi-Objective AHA configuration
  </h2>
  <div class="grid gap-2 grid-cols-1 grid-rows-3">
    <fieldset class="fieldset text-lg">
      <legend class="fieldset-legend">Iterations:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.iterations} />
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg">Population:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.population}/>
    </fieldset>
    <fieldset class="fieldset w-full">
      <legend class="fieldset-legend text-lg text-nowrap">Archive Size:</legend>
      <input type="number" class="input input-lg" placeholder="100" bind:value={config.archiveSize}/>
    </fieldset>
  </div>
</div>
