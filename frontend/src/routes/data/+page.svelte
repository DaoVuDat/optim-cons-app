<script lang="ts">
  import clsx from "clsx";
  import {stepStore} from "$lib/stores/steps.svelte.js";
  import {objectiveStore} from "$lib/stores/objectives.svelte.js";
  import riskConfigComponent from "$lib/components/objective-configs/risk-config.svelte";
  import hoistingConfigComponent from "$lib/components/objective-configs/hoisting-config.svelte";
  import safetyConfigComponent from "$lib/components/objective-configs/safety-config.svelte";
  import {goto} from "$app/navigation";
  import { objectives} from "$lib/wailsjs/go/models";

  const configComponents = {
    [objectives.ObjectiveType.HoistingObjective]: hoistingConfigComponent,
    [objectives.ObjectiveType.RiskObjective]: riskConfigComponent,
    [objectives.ObjectiveType.SafetyObjective]: safetyConfigComponent,
  }

  let selectedObjective = $state<objectives.ObjectiveType>()

  const component = $derived.by(() => {
    if (selectedObjective) {
      return configComponents[selectedObjective]
    }
  })

  let loading = $state<boolean>(false)

  const handleClick = (obj: objectives.ObjectiveType) => {
    selectedObjective = obj;
  }

  const handleNext = async () => {
    loading = true


    // Do loading data objective configs
    await new Promise(() => setTimeout(() => {console.log("process data")}, 2000))
    loading = false

    await goto('/constraint')
    stepStore.nextStep()
  }

  // Hack


</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Data configuration</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 gap-4 w-[1400px] auto-rows-min">
    <div
        class="h-[420px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      {#each objectiveStore.objectives.selectedObjectives as obj (obj)}
        <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer text-left",
        selectedObjective === obj.objectiveType ? 'bg-[#422AD5] text-white' : '')}
                onclick={() => handleClick(obj.objectiveType)}>
          {obj.objectiveType}
        </button>
      {/each}
    </div>
    <div class="h-[420px] overflow-y-auto card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
      {#if selectedObjective}
        {@const Component = component}
        <Component/>
      {:else}
        <p>Please select objectives</p>
      {/if}
    </div>
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="ml-4 btn" href="/problem" onclick={() => stepStore.prevStep()}>Back</a>
    <button class='ml-4 btn'
            onclick={() => handleNext()}
    >Next</button>
  </section>
</div>