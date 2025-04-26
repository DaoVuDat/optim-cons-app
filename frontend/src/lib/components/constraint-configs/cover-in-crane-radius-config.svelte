<script lang="ts">

    import {coverInCraneRadiusConfig} from "$lib/stores/constraints";
    import Modal from "$lib/components/modal.svelte";
    import {type ISelectedCraneWithId} from "$lib/stores/objectives";
    import type {Facility} from "$lib/stores/problems/problem";

    const config = coverInCraneRadiusConfig

    interface Props {
      facilities: Facility[]
    }

    const {facilities}: Props = $props()

    let isOpenModal = $state<boolean>(false)

    let cranes = $state<ISelectedCraneWithId[]>(coverInCraneRadiusConfig.CraneLocations)

    const addCrane = () => {
      cranes.push({
        Id: Math.random().toString(),
        Name: "",
        BuildingNames: "",
        Radius: 0,
        HoistingTimeFilePath: "",
      })
    }

    const removeCrane = (idx: string) => {
      const indexToRemove = cranes.findIndex(crane => crane.Id === idx);
      if (indexToRemove !== -1) {
        cranes.splice(indexToRemove, 1);
      }
    }
</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-1 grid-rows-2 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Power Difference (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="1" bind:value={config.PowerDifferencePenalty}/>
    </fieldset>

    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="20000" bind:value={config.AlphaCoverInCraneRadiusPenalty}/>
    </fieldset>

    <Modal bind:isModalOpen={isOpenModal} buttonText="Save">
      {#snippet content()}
        <div class="h-[600px] overflow-y-auto">
          <div class="grid grid-cols-3 gap-4">
            {#if cranes.length > 0}
              <!--  List of cranes  -->
              {#each cranes as crane, idx (crane.Id)}
                <div class="p-2 card bg-base-100 border shadow-sm flex items-center justify-center">
                  <div class="relative card-body">
                    <div class="absolute text-lg font-bold top-0 left-1">
                      {idx + 1}.
                    </div>
                    <div class="absolute top-0 right-0 card-actions justify-end">
                      <button onclick={() => removeCrane(crane.Id)} aria-label="delete btn"
                              class="btn btn-square btn-xs">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-6 w-6"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor">
                          <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                      </button>
                    </div>

                    <div class="flex flex-col">
                      <div class="flex">
                        <fieldset class="fieldset flex flex-col">
                          <legend class="fieldset-legend text-base">Select crane:</legend>
                          <select class="select select-sm" bind:value={crane.Name}>
                            <option disabled selected></option>
                            {#each facilities as fac (fac)}
                              <option value={fac.Symbol}>{fac.Symbol} - {fac.Name}</option>
                            {/each}
                          </select>
                        </fieldset>
                      </div>
                      <div class="flex items-center">
                        <fieldset class="fieldset flex flex-col">
                          <legend class="fieldset-legend text-base">Facilities in the crane:</legend>
                          <input type="text" class="input input-sm" placeholder="TF1 TF2"
                                 bind:value={crane.BuildingNames}/>
                        </fieldset>
                      </div>
                      <div class="flex items-center">
                        <fieldset class="fieldset flex flex-col">
                          <legend class="fieldset-legend text-base ">Crane's radius:</legend>
                          <input type="number" class="input input-sm" placeholder="40"
                                 bind:value={crane.Radius}/>
                        </fieldset>
                      </div>
                    </div>
                  </div>
                </div>
              {/each}
            {:else}
              <!-- No cranes -->
              <div class="col-start-2 flex items-center justify-center">
                No Crane
              </div>
            {/if}
          </div>
          <div class="mt-8 flex justify-center items-center">
            <button onclick={addCrane} class="btn btn-soft btn-primary">Add Crane</button>
          </div>
        </div>
      {/snippet}
    </Modal>
  </div>
  <div class="flex justify-end items-center">
    <button class="btn btn-primary" onclick={()=>isOpenModal = true}>Setup Cranes</button>
  </div>
</div>