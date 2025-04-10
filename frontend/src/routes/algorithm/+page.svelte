<script lang="ts">
  import clsx from "clsx";
  import moahaConfig from '$lib/components/algo-configs/moaha-config.svelte'
  import gaConfig from '$lib/components/algo-configs/ga-config.svelte'
  import ahaConfig from '$lib/components/algo-configs/aha-config.svelte'
  import gwoConfig from '$lib/components/algo-configs/gwo-config.svelte'
  import {stepStore} from "$lib/stores/steps.svelte.js";
  import {algorithmsStore, type AlgorithmWithLabel} from "$lib/stores/algorithms.svelte";
  import {algorithms, main} from "$lib/wailsjs/go/models";
  import {goto} from "$app/navigation";
  import {ConstraintsInfo, ProblemInfo} from "$lib/wailsjs/go/main/App";

  const configComponents = {
    [algorithms.AlgorithmType.MOAHA]: moahaConfig,
    [algorithms.AlgorithmType.AHA]: ahaConfig,
    [algorithms.AlgorithmType.GWO]: gwoConfig,
    [algorithms.AlgorithmType.GeneticAlgorithm]: gaConfig,
  }

  const component = $derived.by(() => {
    if (algorithmsStore.getValidSelection()) {
      return configComponents[algorithmsStore.selectedAlgorithm!.value]
    }
  })

  const handleClick = (algo: AlgorithmWithLabel) => {
    algorithmsStore.selectedAlgorithm = algo;
  }

  let loading = $state<boolean>(false)

  const handleNext = async () => {
    loading = true

    const problemInfo = await ProblemInfo()
    const constraintsInfo = await ConstraintsInfo()
    console.log(problemInfo)
    console.log(constraintsInfo)

    // Create Algo
    // await AddConstraints(constraintsInput)

    loading = false

    await goto('/optimize')

    stepStore.nextStep()
  }

</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Select algorithm</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 gap-4 w-[1400px] auto-rows-min">
    <div
        class="h-[420px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      {#each algorithmsStore.validAlgorithmsList as algo (algo)}
        <button class={clsx("p-4 rounded h-18 flex justify-between items-center cursor-pointer text-left",
        algorithmsStore.selectedAlgorithm?.value === algo.value ? 'bg-[#422AD5] text-white' : '')}
                onclick={() => handleClick(algo)}>
          {algo.label}
        </button>
      {/each}
    </div>
    <div class="h-[420px] card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
      {#if algorithmsStore.getValidSelection()}
        {@const Component = component}
        <Component/>
      {:else}
        <p>Please select an algorithm</p>
      {/if}
    </div>

    {#if loading}
      <div class="toast toast-center toast-middle">
        <div class="alert alert-info">
          <span>Setting up algorithm...</span>
        </div>
      </div>
    {/if}
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="btn" href="/constraint" onclick={() => stepStore.prevStep()}>Back</a>
    <button class={clsx('ml-4 btn', algorithmsStore.getValidSelection() ? '' : 'btn-disabled')}
       onclick={handleNext}
    >Next</button>
  </section>
</div>