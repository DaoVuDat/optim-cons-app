<script>
import {stepStore} from "$lib/stores/steps.svelte.js";
import clsx from "clsx";
import {problemStore} from "$lib/stores/problem.svelte.js";
</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Add constraints</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 grid-rows-2 gap-4 w-[1600px] auto-rows-min">
    <div
        class="h-[280px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">

    </div>
    <div class="card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">

    </div>

    <div
        class="h-[280px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">

    </div>
    <div class="card p-4 bg-base-100 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">

    </div>
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class="btn" href="/data" onclick={() => stepStore.prevStep()}>Back</a>
    <a class={clsx('ml-4 btn', true ? '' : 'btn-disabled')}
       href="/algorithm"
       onclick={() => stepStore.nextStep()}
    >Next</a>
  </section>
</div>