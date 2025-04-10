<script lang="ts">
  import {stepStore} from "$lib/stores/steps.svelte";
  import {AlgorithmInfo, RunAlgorithm} from "$lib/wailsjs/go/main/App";
  import {onDestroy, onMount} from "svelte";
  import {EventsOff, EventsOn} from "$lib/wailsjs/runtime";
  import {main} from "$lib/wailsjs/go/models";
  import {algorithmsStore} from "$lib/stores/algorithms.svelte";
  import clsx from "clsx";
  import type {Coordinate, Facility} from "$lib/stores/problems/problem";

  let progress = $state<number>(0)
  let progressInfo = $state<string>("")
  let isMulti = $state<boolean>(false)

  const handleOptimize = async () => {
    const algorithmInfo = await AlgorithmInfo()
    console.log(algorithmInfo)

    await RunAlgorithm()
  }

  interface Progress {
    progress: number
  }

  type MultiObjective = {
    numberOfAgentsInArchive: number
    type: 'multi'
  } & Progress

  type SingleObjective = {
    bestFitness: number
    type: 'single'

  } & Progress

  const roundNDecimal = (value: number, n: number) => {
    const pow = Math.pow(10, n)
    return Math.round((value + Number.EPSILON) * pow) / pow
  }

  onMount(() => {
    // Listen for the 'backendEvent' emitted from Go
    EventsOn(main.EventType.ProgressEvent, (data: MultiObjective | SingleObjective) => {
      if (data) {
        // check the type of problem ( single or multiple )
        if (data.type === 'multi') {
          isMulti = true;
          progress = Math.round(data.progress)
          progressInfo = `${data.numberOfAgentsInArchive}`
        } else if (data.type === 'single') {
          isMulti = false;
          progress = Math.round(data.progress)
          progressInfo = `${roundNDecimal(data.bestFitness, 4)}`
        }
      }
    });

    EventsOn(main.EventType.ResultEvent, (data: ResultLocation[]) => {
      if (data) {
        results.push(...data.map(r => ({
          ...r,
          Id: Math.random()
        })))

        console.log(results.length)
      }
    });
  })

  onDestroy(() => {
    EventsOff(main.EventType.ProgressEvent)
    EventsOff(main.EventType.ResultEvent)
  })

  interface Penalty {
    [k: string]: number
  }

  interface MapLocation {
    [k: string]: Facility
  }

  interface ResultLocation {
    MapLocations: MapLocation
    Value: number[]
    Penalty: Penalty
  }

  interface ResultLocationWithId extends ResultLocation {
    Id: number
  }

  let results = $state<ResultLocationWithId[]>([])
  let selectedResult = $state<ResultLocationWithId|undefined>(undefined)

  const handleSelectedResult = (result: ResultLocationWithId) => {
    selectedResult = result
  }

  $inspect(selectedResult, results)
</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Optimize</h1>
  </section>


  <!-- Content -->
  <section class="h-[480px] px-24 grid grid-cols-12 grid-rows-3 gap-4 w-[1400px] auto-rows-min ">
    <div class="pl-2 py-4 row-start-1 col-start-1 col-span-4 card bg-base-100 shadow-md rounded-lg flex flex-col justify-center items-center">
      <div class="w-full px-4 flex items-center justify-center mb-2">
        <progress class="progress progress-info w-full" value={progress} max="100">
        </progress>
        <div class="pl-4 pr-2 w-16">{progress}%</div>
      </div>

      {#if isMulti}
        <div class="flex flex-col justify-center items-center">
          <span>Number of solutions:</span>
          <span>{progressInfo} &nbsp;</span>
        </div>
      {:else}
        <div class="flex flex-col justify-center items-center">
          <span>Best result (minimum):</span>
          <span>{progressInfo} &nbsp;</span>
        </div>
      {/if}

    </div>
    <div
        class="px-2 py-4 col-start-5 row-start-1 col-span-8 row-span-3 card bg-base-100 shadow-md rounded-lg flex flex-col justify-center items-center">

    </div>
    <div
        class="px-2 py-4 max-h-full col-start-1 row-start-2 row-span-2 col-span-4 card bg-base-100 shadow-md rounded-lg flex flex-col overflow-y-auto">
      {#each results as res, idx (res.Id + idx)}
        <button class={clsx("p-4 rounded h-18 flex justify-between items-center cursor-pointer text-left",
      selectedResult?.Id === res.Id ? 'bg-[#422AD5] text-white' : '')}
                onclick={() => handleSelectedResult(res)}>
          Result #{idx+1} ({Object.values(res.Penalty).reduce((prev, cur) => prev+cur, 0) !== 0 ?  "Infeasible" : "Feasible"})
        </button>
      {/each}
    </div>
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="btn" href="/algorithm" onclick={() => stepStore.prevStep()}>Back</a>
    <button class='ml-4 btn' onclick={handleOptimize}>Optimize</button>
  </section>
</div>