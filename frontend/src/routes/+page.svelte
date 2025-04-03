<script lang="ts">
  import clsx from "clsx";
  import {objectiveStore, ObjectiveType} from "$lib/stores/objectives.svelte";
  import {SelectFile} from "$lib/wailsjs/go/main/App";

  let disable = $state(false)

  let data = $state("")

  const selectFile = async () => {
    disable = !disable

    data = await SelectFile()

    disable = !disable
  }

  $inspect(objectiveStore.objectiveList, objectiveStore.selectObjectiveOptions,
  objectiveStore.selectedObjectiveList,objectiveStore.selectedObjectiveOptions)



</script>
<div class="h-[calc(100vh-64px-64px)] w-full pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Select objectives</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-13 gap-4 w-[900px] auto-rows-min">
    <div class="h-96 shadow rounded col-span-6">
      <select class="flex w-full h-full text-black p-4 borderless" multiple bind:value={objectiveStore.selectObjectiveOptions}>
        {#each objectiveStore.objectiveList as s (s.value)}
          <option class="styled-option" value={s.value}>{s.label}</option>
        {/each}
      </select>
    </div>
    <div class="flex justify-center items-center col-span-1">
      <button class="btn" onclick={() => console.log(objectiveStore.selectObjectiveOptions)}>&rarr;</button>
    </div>
    <div class="h-96 shadow rounded col-span-6">
      <select class="w-full h-full text-black p-4 borderless" multiple bind:value={objectiveStore.selectedObjectiveOptions}>
        {#each objectiveStore.selectedObjectiveList as s (s.value)}
          <option  value={s.value}>{s.value}</option>
        {/each}
      </select>
    </div>
    <div class="h-48 bg-red-400 col-span-full">Info</div>
    <!--    <button class={clsx("btn", {-->
    <!--  "btn-disabled": disable-->
    <!--})}-->
    <!--            onclick={selectFile}-->
    <!--    >-->
    <!--      <span class={disable ? "loading loading-ring loading-sm" : ""}></span>-->
    <!--      Select file-->
    <!--    </button>-->
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class={clsx('btn', objectiveStore.objectives.numberOfObjectives.length === 0 ? 'btn-disabled': '')}
       href="/algorithm">Next</a>
  </section>
</div>


<style>
  .borderless {
      border: none;
      outline: none;
      appearance: none;
      -webkit-appearance: none;
      -moz-appearance: none;
      background-color: transparent;
  }

  .styled-option {
      height: 50px;
      padding: 10px;
      margin: 2px 0;
      line-height: 1.5;
      border-bottom: 1px solid #eee;
      color: #333;
  }

</style>