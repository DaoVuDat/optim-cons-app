<script lang="ts">
  import clsx from "clsx";
  import moahaConfig from '$lib/components/moaha-config.svelte'
  import gaConfig from '$lib/components/ga-config.svelte'
  import ahaConfig from '$lib/components/aha-config.svelte'
  import gwoConfig from '$lib/components/gwo-config.svelte'
  import {stepStore} from "$lib/stores/steps.svelte.js";
  import {Algorithms, algorithmsStore, type AlgorithmWithLabel} from "$lib/stores/algorithms.svelte";

  const configComponents = {
    [Algorithms.MOAHA]: moahaConfig,
    [Algorithms.AHA]: ahaConfig,
    [Algorithms.GWO]: gwoConfig,
    [Algorithms.GA]: gaConfig,
  }

  const component = $derived.by(() => {
    if (algorithmsStore.getValidSelection()) {
      return configComponents[algorithmsStore.selectedAlgorithm!.value]
    }
  })

  const handleClick = (algo: AlgorithmWithLabel) => {
    algorithmsStore.selectedAlgorithm = algo;
  }

  $inspect(algorithmsStore.selectedAlgorithm)

</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Select algorithm</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 gap-4 w-[1600px] auto-rows-min">
    <div
        class="h-[580px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      {#each algorithmsStore.validAlgorithmsList as algo (algo)}
        <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer text-left",
        algorithmsStore.selectedAlgorithm?.value === algo.value ? 'bg-[#422AD5] text-white' : '')}
                onclick={() => handleClick(algo)}>
          {algo.label}
        </button>
      {/each}
    </div>
    <div class="card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
      {#if algorithmsStore.getValidSelection()}
        {@const Component = component}
        <Component/>
      {:else}
        <p>Please select an algorithm</p>
      {/if}
    </div>
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="btn" href="/" onclick={() => stepStore.prevStep()}>Back</a>
    <a class={clsx('ml-4 btn', algorithmsStore.getValidSelection() ? '' : 'btn-disabled')}
       href="/problem"
       onclick={() => stepStore.nextStep()}
    >Next</a>
  </section>
</div>