<script lang="ts">
  import clsx from "clsx";
  import {type IOptions, objectiveStore} from "$lib/stores/objectives.svelte";
  import {stepStore} from "$lib/stores/steps.svelte";
  import {SelectFile} from "$lib/wailsjs/go/main/App";

  let disable = $state(false)

  let data = $state("")

  const selectFile = async () => {
    disable = !disable

    data = await SelectFile()

    disable = !disable
  }

  const handleClick = (option: IOptions) => {
    objectiveStore.selectObjectiveOption = option
  }

</script>
<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
  <!-- Top Section -->
  <section class="mt-8 text-black">
    <h1 class="text-5xl font-bold">Select objectives</h1>
  </section>


  <!-- Content -->
  <section class="px-24 grid grid-cols-12 gap-4 w-[1600px] auto-rows-min">
    <div class="h-[580px] px-2 py-4 card bg-base-100 shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
      {#each objectiveStore.objectiveList as s (s.value)}
        <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer",
          s.value === objectiveStore.selectObjectiveOption?.value ? 'bg-[#422AD5] text-white' : ''
        )}
                onclick={() => handleClick(s)}
        >
          {s.label}
          <input type="checkbox" class="custom-checkbox" bind:checked={s.isChecked} onchange={() => {
            objectiveStore.selectObjective(s)
          }}/>
        </button>
      {/each}
    </div>
    <div class="card p-4 bg-base-100 shadow-md rounded-lg col-span-8">{
        objectiveStore.selectObjectiveOption?.content
    }</div>
  </section>

  <!-- Bottom Section -->
  <section class="w-full text-end">
    <a class={clsx('btn', objectiveStore.objectives.selectedObjectives.length === 0 ? 'btn-disabled': '')}
       href="/algorithm" onclick={() => stepStore.nextStep()}>Next</a>
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