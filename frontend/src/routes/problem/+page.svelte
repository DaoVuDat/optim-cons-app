<script lang="ts">
  import clsx from "clsx";
  import {stepStore} from "$lib/stores/steps.svelte.js";
  import {problemList, problemStore, type ProblemWithLabel} from "$lib/stores/problem.svelte.js";
  import continuousProblemConfigComponent from "$lib/components/problem-configs/continuous-config.svelte";
  import gridProblemConfigComponent from "$lib/components/problem-configs/grid-config.svelte";
  import PredeterminedConfig from "$lib/components/problem-configs/predetermined-config.svelte";
  import {goto} from "$app/navigation";
  import {CreateProblem} from "$lib/wailsjs/go/main/App";
  import {main, data as dataType, conslay_predetermined} from "$lib/wailsjs/go/models";
  import {predeterminedProblemConfig} from "$lib/stores/problems";
  import {toast} from "@zerodevx/svelte-toast";
  import {errorOpts, infoOpts, successOpts} from "$lib/utils/toast-opts";


  const configComponents = {
    [dataType.ProblemName.ContinuousConstructionLayout]: continuousProblemConfigComponent,
    [dataType.ProblemName.GridConstructionLayout]: gridProblemConfigComponent,
    [dataType.ProblemName.PredeterminedConstructionLayout]: PredeterminedConfig,
  }

  const component = $derived.by(() => {
    if (problemStore.getValidSelection()) {
      return configComponents[problemStore.selectedProblem!.value]
    }
  })

  let loading = $state<boolean>(false)

  const handleClick = (prob: ProblemWithLabel) => {
    problemStore.selectedProblem = prob;
  }

  const handleNext = async () => {
    loading = true
    toast.push("Setting up problem...", {
      theme: infoOpts
    })
    try {
      if (problemStore.selectedProblem) {
        switch (problemStore.selectedProblem!.value) {
          case dataType.ProblemName.ContinuousConstructionLayout : {
            const config = problemStore.getConfig(problemStore.selectedProblem.value)

            const problemInput = new main.ProblemInput({
              problemName: problemStore.selectedProblem!.value,
              layoutLength: config.length,
              layoutWidth: config.width,
              facilitiesFilePath: config.facilitiesFilePath.value,
              phasesFilePath: config.phasesFilePath.value
            })
            await CreateProblem(problemInput)
            break
          }
          case dataType.ProblemName.GridConstructionLayout : {
            const config = problemStore.getConfig(problemStore.selectedProblem.value)

            const problemInput = new main.ProblemInput({
              problemName: problemStore.selectedProblem!.value,
              layoutLength: config.length,
              layoutWidth: config.width,
              facilitiesFilePath: config.facilitiesFilePath.value,
              phasesFilePath: config.phasesFilePath.value,
              gridSize: config.gridSize,
            })
            await CreateProblem(problemInput)
            break
          }
          case dataType.ProblemName.PredeterminedConstructionLayout:

            const config = predeterminedProblemConfig
            const fixedFacilities: conslay_predetermined.LocFac[] = config.fixedFacilities.filter(f => f.FacilityName).map((f) => ({
              facName: f.FacilityName,
              locName: f.LocName,
            }))


            const problemInput = new main.ProblemInput({
              problemName: problemStore.selectedProblem!.value,
              numberOfLocations: config.value.numberOfLocations,
              numberOfFacilities: config.value.numberOfFacilities,
              fixedFacilities: fixedFacilities,
            })
            await CreateProblem(problemInput)
            break
        }


        await goto('/data')
        stepStore.nextStep()
        toast.push("Set up problem!", {
          theme: successOpts
        })
      }
    } catch (err) {
      toast.pop(0)
      toast.push(err as string, {
        theme: errorOpts
      })
    } finally {
      loading = false
    }


  }


</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <!--    <section class="mt-8 text-black">-->
  <!--        <h1 class="text-5xl font-bold">Select problem</h1>-->
  <!--    </section>-->

  <!-- Content -->
  <section class="mt-8 px-24 grid grid-cols-12 gap-4 w-[1400px] auto-rows-min">
    <div
        class="h-[560px] bg-base-100 px-2 py-4 card shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      {#each problemStore.validProblemList as prob (prob)}
        <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer text-left",
        problemStore.selectedProblem?.value === prob.value ? 'bg-[#422AD5] text-white' : '')}
                onclick={() => handleClick(prob)}>
          {prob.label}
        </button>
      {/each}
    </div>
    <div
        class="h-[560px] bg-base-100 overflow-y-auto card p-4 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
      {#if problemStore.getValidSelection()}
        {@const Component = component}
        <Component/>
      {:else}
        <p>Please select problem</p>
      {/if}
    </div>

    {#if loading}
      <div class="toast toast-center toast-middle">
        <div class="alert alert-info">
          <span>Setting up problem...</span>
        </div>
      </div>
    {/if}
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="ml-4 btn" href="/" onclick={() => {
        stepStore.prevStep()
        problemStore.resetSelection()
    }}>Back</a>
    <button class={clsx('ml-4 btn', problemStore.getValidSelection() ? '' : 'btn-disabled')}
            onclick={() => handleNext()}
    >Next
    </button>
  </section>
</div>