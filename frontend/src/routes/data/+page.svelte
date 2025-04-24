<script lang="ts">
  import clsx from "clsx";
  import {stepStore} from "$lib/stores/steps.svelte.js";
  import {objectiveStore} from "$lib/stores/objectives.svelte.js";
  import riskConfigComponent from "$lib/components/objective-configs/risk-config.svelte";
  import hoistingConfigComponent from "$lib/components/objective-configs/hoisting-config.svelte";
  import safetyConfigComponent from "$lib/components/objective-configs/safety-config.svelte";
  import safetyHazardConfigComponent from "$lib/components/objective-configs/safety-hazard-config.svelte";
  import transportCostConfigComponent from "$lib/components/objective-configs/transport-cost-config.svelte";
  import constructionCostConfigComponent from "$lib/components/objective-configs/construction-cost-config.svelte";
  import {goto} from "$app/navigation";
  import {main, data as dataType} from "$lib/wailsjs/go/models";
  import type {PageProps} from "../../../.svelte-kit/types/src/routes/data/$types";
  import type {Facility} from "$lib/stores/problems/problem";
  import {CreateObjectives} from "$lib/wailsjs/go/main/App";


  const configComponents = {
    [dataType.ObjectiveType.HoistingObjective]: hoistingConfigComponent,
    [dataType.ObjectiveType.RiskObjective]: riskConfigComponent,
    [dataType.ObjectiveType.SafetyObjective]: safetyConfigComponent,
    [dataType.ObjectiveType.TransportCostObjective]: transportCostConfigComponent,
    [dataType.ObjectiveType.SafetyHazardObjective]: safetyHazardConfigComponent,
    [dataType.ObjectiveType.ConstructionCostObjective]: constructionCostConfigComponent
  }

  let selectedObjective = $state<dataType.ObjectiveType>()

  const component = $derived.by(() => {
    if (selectedObjective) {
      return configComponents[selectedObjective]
    }
  })

  let loading = $state<boolean>(false)

  const handleClick = (obj: dataType.ObjectiveType) => {
    selectedObjective = obj;
  }

  const handleNext = async () => {
    loading = true

    // Do loading data objective configs
    if (objectiveStore.objectives.selectedObjectives.length > 0) {
      const objectivesInput = objectiveStore.objectives.selectedObjectives.map<main.ObjectiveInput>(objective => {
        return {
          objectiveName: objective.objectiveType,
          objectiveConfig: objective.config,
        }
      })

      await CreateObjectives(objectivesInput)
    }
    loading = false

    await goto('/constraint')
    stepStore.nextStep()
  }

  let {data}: PageProps = $props();

  let facilities = $state<Facility[] | undefined>(undefined)

  if (data.problemInfo.problemName === dataType.ProblemName.ContinuousConstructionLayout ||
    data.problemInfo.problemName === dataType.ProblemName.GridConstructionLayout) {
    // convert map to array
    facilities = Object.values(data.problemInfo.locations)
  }

  // Hack + make config "reactive"
  const noTypeCheck = (x: any) => {
    if (x) {
      return x
    } else
      return {}
  };

</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
<!--  <section class="mt-8 text-black">-->
<!--    <h1 class="text-5xl font-bold">Data configuration</h1>-->
<!--  </section>-->


  <!-- Content -->
  <section class="mt-8 px-24 grid grid-cols-12 gap-4 w-[1400px] auto-rows-min">
    <div
        class="h-[560px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      {#each objectiveStore.objectives.selectedObjectives as obj (obj)}
        <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer text-left",
        selectedObjective === obj.objectiveType ? 'bg-[#422AD5] text-white' : '')}
                onclick={() => handleClick(obj.objectiveType)}>
          {obj.objectiveType}
        </button>
      {/each}
    </div>
    <div
        class="h-[560px] overflow-y-auto card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
      {#if selectedObjective}
        {@const Component = component}
        <Component {...noTypeCheck({
          facilities
        })}/>
      {:else}
        <p>Please select objectives</p>
      {/if}
    </div>

    {#if loading}
      <div class="toast toast-center toast-middle">
        <div class="alert alert-info">
          <span>Loading data...</span>
        </div>
      </div>
    {/if}

  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="ml-4 btn" href="/problem" onclick={() => stepStore.prevStep()}>Back</a>
    <button class='ml-4 btn'
            onclick={() => handleNext()}
    >Next
    </button>
  </section>
</div>