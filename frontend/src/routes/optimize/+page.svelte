<script lang="ts">
import {stepStore} from "$lib/stores/steps.svelte";
import {AlgorithmInfo, RunAlgorithm} from "$lib/wailsjs/go/main/App";
import {onMount} from "svelte";
import {EventsOn} from "$lib/wailsjs/runtime";
import {main} from "$lib/wailsjs/go/models";

let progress = $state<number>(0)

const handleOptimize = async () => {
  const algorithmInfo = await AlgorithmInfo()
  console.log(algorithmInfo)

  await RunAlgorithm()
}

onMount(() => {
  // Listen for the 'backendEvent' emitted from Go

  EventsOn("ProgressEvent", (data) => {
    if (data) {
      console.log('Received event from backend:', data);
    }
  });
})

</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Optimize</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 gap-4 w-[1400px] auto-rows-min">
    <div>
      <div class="radial-progress" style="--value:{progress}; --size:12rem; --thickness: 2rem;"
           aria-valuenow={progress}
           role="progressbar">
        {progress}%
      </div>
    </div>
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="btn" href="/constraint" onclick={() => stepStore.prevStep()}>Back</a>
    <button class='ml-4 btn' onclick={handleOptimize}>Optimize</button>
  </section>
</div>