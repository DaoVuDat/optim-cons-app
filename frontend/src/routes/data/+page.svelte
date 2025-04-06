<script lang="ts">
  import clsx from "clsx";
  import {stepStore} from "$lib/stores/steps.svelte.js";
  // import {problemList, problemStore, ProblemType, type ProblemWithLabel} from "$lib/stores/problem.svelte.js";
  import continuousProblemConfigComponent from "$lib/components/problem-configs/continuous-config.svelte";
  import gridProblemConfigComponent from "$lib/components/problem-configs/grid-config.svelte";
  import PredeterminatedConfig from "$lib/components/problem-configs/predeterminated-config.svelte";
  import {continuousProblemConfig} from "$lib/stores/problems";
  import {goto} from "$app/navigation";

  // const configComponents = {
  //   [ProblemType.Continuous]: continuousProblemConfigComponent,
  //   [ProblemType.Grid]: gridProblemConfigComponent,
  //   [ProblemType.PreLocated]: predeterminedProblemConfigComponent,
  // }
  //
  // const component = $derived.by(() => {
  //   if (problemStore.getValidSelection()) {
  //     return configComponents[problemStore.selectedProblem!.value]
  //   }
  // })
  //
  let loading = $state<boolean>(false)
  //
  // const handleClick = (prob: ProblemWithLabel) => {
  //   problemStore.selectedProblem = prob;
  // }

  const handleNext = async () => {
    loading = true

    await new Promise(() => setTimeout(() => {console.log("process data")}, 2000))
    loading = false

    await goto('/data')
    stepStore.nextStep()
  }

</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Data configuration</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 gap-4 w-[1600px] auto-rows-min">
    <div
        class="h-[580px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      <!--{#each problemList as prob (prob)}-->
      <!--  <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer text-left",-->
      <!--  problemStore.selectedProblem?.value === prob.value ? 'bg-[#422AD5] text-white' : '')}-->
      <!--          onclick={() => handleClick(prob)}>-->
      <!--    {prob.label}-->
      <!--  </button>-->
      <!--{/each}-->
    </div>
    <div class="h-[580px] overflow-y-auto card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
      <!--{#if problemStore.getValidSelection()}-->
      <!--  {@const Component = component}-->
      <!--  <Component/>-->
      <!--{:else}-->
      <!--  <p>Please select problem</p>-->
      <!--{/if}-->
    </div>
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="ml-4 btn" href="/problem" onclick={() => stepStore.prevStep()}>Back</a>
<!--    <button class={clsx('ml-4 btn', problemStore.getValidSelection() ? '' : 'btn-disabled')}-->
<!--            onclick={() => handleNext()}-->
<!--    >Next</button>-->
  </section>
</div>