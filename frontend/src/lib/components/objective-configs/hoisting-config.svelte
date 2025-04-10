<script lang="ts">
    import {SelectFile} from "$lib/wailsjs/go/main/App";
    import {hoistingConfig, type ISelectedCrane, type ISelectedCraneWithId} from "$lib/stores/objectives";
    import Modal from "$lib/components/modal.svelte"
    import type {Facility} from "$lib/stores/problems/problem";

    interface Props {
        facilities: Facility[]
    }

    const {facilities}: Props = $props()

    let cranes = $state<ISelectedCraneWithId[]>(hoistingConfig.CraneLocations)

    const config = hoistingConfig

    let isOpenModal = $state<boolean>(false)

    const selectFile = async (idx: string) => {
        cranes.find(crane => crane.Id === idx)!.HoistingTimeFilePath = await SelectFile()
    }

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
        cranes = cranes.filter(crane => crane.Id !== idx)
    }

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-4 grid-rows-3 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Number of Floors:</legend>
      <input type="number" class="input input-lg" placeholder="10" bind:value={config.NumberOfFloors}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Floor height:</legend>
      <input type="number" class="input input-lg" placeholder="3.2" bind:value={config.FloorHeight}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">ZM:</legend>
      <input type="number" class="input input-lg" placeholder="2" bind:value={config.ZM}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Vuvg:</legend>
      <input type="number" class="input input-lg" placeholder="37.5" bind:value={config.Vuvg}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Vlvg:</legend>
      <input type="number" class="input input-lg" placeholder="18.75" bind:value={config.Vlvg}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Vag:</legend>
      <input type="number" class="input input-lg" placeholder="50" bind:value={config.Vag}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Vwg:</legend>
      <input type="number" class="input input-lg" placeholder="0.5" bind:value={config.Vwg}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha Hoisting:</legend>
      <input type="number" class="input input-lg" placeholder="0.25" bind:value={config.AlphaHoisting}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Beta Hoisting:</legend>
      <input type="number" class="input input-lg" placeholder="1" bind:value={config.BetaHoisting}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">NHoisting:</legend>
      <input type="number" class="input input-lg" placeholder="1" bind:value={config.NHoisting}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="1" bind:value={config.AlphaHoistingPenalty}/>
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
                      <div class="flex items-center">
                        <fieldset class="fieldset flex flex-col">
                          <legend class="fieldset-legend text-base">Hoisting time file:</legend>
                          <div class="join">
                            <div>
                              <label class="input input-sm validator join-item">
                                <input type="text" placeholder="path://" bind:value={crane.HoistingTimeFilePath}/>
                              </label>
                            </div>
                            <button class="btn btn-neutral join-item btn-sm"
                                    onclick={() =>selectFile(crane.Id)}>Select file
                            </button>
                          </div>
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
      {#snippet moreButtons()}
        <button class="btn  btn-primary">Import Data Template</button>
      {/snippet}
    </Modal>
  </div>
  <div class="flex justify-end items-center">
    <button class="btn btn-primary" onclick={()=>isOpenModal = true}>Setup Cranes</button>
  </div>
</div>