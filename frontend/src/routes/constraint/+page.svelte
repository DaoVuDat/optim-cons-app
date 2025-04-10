<script lang="ts">
  import clsx from "clsx";
  import {stepStore} from "$lib/stores/steps.svelte.js";
  import {goto} from "$app/navigation";
  import {main, data as dataType} from "$lib/wailsjs/go/models";
  import type {PageProps} from "../../../.svelte-kit/types/src/routes/data/$types";
  import type {Facility} from "$lib/stores/problems/problem";
  import outOfBoundConfigComponent from "$lib/components/constraint-configs/out-of-bound-config.svelte"
  import overlapConfigComponent  from "$lib/components/constraint-configs/overlap-config.svelte"
  import coverInCraneRadiusConfigComponent from "$lib/components/constraint-configs/cover-in-crane-radius-config.svelte"
  import inclusiveZoneConfigComponent from "$lib/components/constraint-configs/inclusive-zone-config.svelte"
  import {
    AddConstraints,
  } from "$lib/wailsjs/go/main/App";

  import {constraintsStore} from "$lib/stores/constraints.svelte";

  const configComponents = {
    [dataType.ConstraintType.OutOfBound]: outOfBoundConfigComponent,
    [dataType.ConstraintType.Overlap]: overlapConfigComponent,
    [dataType.ConstraintType.InclusiveZone]: inclusiveZoneConfigComponent,
    [dataType.ConstraintType.CoverInCraneRadius]: coverInCraneRadiusConfigComponent,
  }

  let selectedConstraint = $state<dataType.ConstraintType>()

  const component = $derived.by(() => {
    if (selectedConstraint) {
      return configComponents[selectedConstraint]
    }
  })

  let loading = $state<boolean>(false)

  const handleClick = (obj: dataType.ConstraintType) => {
    selectedConstraint = obj;
  }

  const handleNext = async () => {
    loading = true

    console.log($state.snapshot(constraintsStore.constraints))

    // Do loading data objective configs
    if (constraintsStore.constraints.selectedConstraints.length > 0) {
      const constraintsInput = constraintsStore.constraints.selectedConstraints.map<main.ConstraintInput>(con => {
        return {
          constraintName: con.constraintType,
          constraintConfig: con.config,
        }
      })

      await AddConstraints(constraintsInput)
    }

    loading = false

    await goto('/algorithm')
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
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Data configuration</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 gap-4 w-[1400px] auto-rows-min">
    <div
        class="h-[420px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      {#each constraintsStore.constraintList as con (con)}
        <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer",
          con.value === selectedConstraint ? 'bg-[#422AD5] text-white' : ''
        )}
                onclick={() => handleClick(con.value)}
        >
          {con.label}
          <input type="checkbox" class="custom-checkbox" bind:checked={con.isChecked} onchange={() => {
            constraintsStore.selectConstraint(con)
          }}/>
        </button>
      {/each}
    </div>
    <div
        class="h-[420px] overflow-y-auto card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
      {#if selectedConstraint}
        {@const Component = component}
        <Component {...noTypeCheck({
          facilities
        })}/>
      {:else}
        <p>Please select constraint</p>
      {/if}
    </div>

    {#if loading}
      <div class="toast toast-center toast-middle">
        <div class="alert alert-info">
          <span>Adding constraints...</span>
        </div>
      </div>
    {/if}
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="ml-4 btn" href="/data" onclick={() => stepStore.prevStep()}>Back</a>
    <button class='ml-4 btn'
            onclick={() => handleNext()}
    >Next
    </button>
  </section>
</div>

<style>
    .custom-checkbox {
        width: 20px;
        height: 20px;
        cursor: pointer;
    }

    /* Optional: more custom tick styling (for full control) */
    /* This part is only needed if you want to fully customize the tick */
    .custom-checkbox {
        appearance: none;
        border: 2px solid #999;
        border-radius: 5px;
        background-color: white;
        position: relative;
    }

    .custom-checkbox:checked {
        background-color: white;
        border-color: black;
    }

    .custom-checkbox:checked::after {
        content: '';
        position: absolute;
        top: 50%;
        left: 50%;
        width: 6px;
        height: 12px;
        border: solid #422AD5;
        border-width: 0 2px 2px 0;
        transform: translate(-48%, -59%) rotate(45deg);
    }
</style>