<script lang="ts">
  import {Algorithms, algorithmsStore} from "$lib/stores/algorithms.svelte";

  let config = algorithmsStore.getConfig(Algorithms.MOAHA)

  // let isModalOpen = $state(false)
</script>

<div class="p-2 w-full h-full">
  <h2 class="text-2xl text-center font-bold">
    Multi-Objective AHA configuration
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
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg">Archive Size:</legend>
      <input type="text" class="input input-lg" placeholder="100" bind:value={config.archiveSize.value}/>
    </fieldset>
  </div>
</div>


<!--&lt;!&ndash; Open the modal using ID.showModal() method &ndash;&gt;-->
<!--<button class="btn modal-button" onclick={()=>isModalOpen = true}>open modal</button>-->

<!--<Modal bind:isModalOpen>-->
<!--  {#snippet content()}-->
<!--    <h3 class="font-bold text-lg">Congratulations random Interner user!</h3>-->
<!--    <p class="py-4">You've been selected for a chance to get one year of subscription to use Wikipedia for free!</p>-->
<!--    <h3 class="font-bold text-lg">Congratulations random Interner user!</h3>-->
<!--    <p class="py-4">You've been selected for a chance to get one year of subscription to use Wikipedia for free!</p>-->
<!--    <h3 class="font-bold text-lg">Congratulations random Interner user!</h3>-->
<!--    <p class="py-4">You've been selected for a chance to get one year of subscription to use Wikipedia for free!</p>-->
<!--    <h3 class="font-bold text-lg">Congratulations random Interner user!</h3>-->
<!--    <p class="py-4">You've been selected for a chance to get one year of subscription to use Wikipedia for free!</p><h3 class="font-bold text-lg">Congratulations random Interner user!</h3>-->
<!--    <p class="py-4">You've been selected for a chance to get one year of subscription to use Wikipedia for free!</p>-->
<!--    <h3 class="font-bold text-lg">Congratulations random Interner user!</h3>-->
<!--    <p class="py-4">You've been selected for a chance to get one year of subscription to use Wikipedia for free!</p>-->

<!--  {/snippet}-->
<!--</Modal>-->