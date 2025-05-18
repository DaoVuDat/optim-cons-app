<script lang="ts">
    import {SelectFile} from "$lib/wailsjs/go/main/App";
    import {hoistingConfig, type ISelectedCrane, type ISelectedCraneWithId, type Building} from "$lib/stores/objectives";
    import Modal from "$lib/components/modal.svelte"
    import type {Facility} from "$lib/stores/problems/problem";

    interface Props {
        facilities: Facility[]
    }

    const {facilities}: Props = $props()

    let cranes = $state<ISelectedCraneWithId[]>(hoistingConfig.CraneLocations)
    let buildings = $state<Building[]>(hoistingConfig.Buildings)

    const config = hoistingConfig

    // Sync cranes and buildings with the store
    $effect(() => {
        hoistingConfig.CraneLocations = cranes
        hoistingConfig.Buildings = buildings
    })

    let isOpenModal = $state<boolean>(false)
    let currentModalStep = $state<'buildings' | 'cranes'>('buildings')

    // New building form data
    let newBuilding = $state<Building>({
        NumberOfFloors: 0,
        FloorHeight: 0,
        Name: ""
    })

    const selectFile = async (idx: string) => {
        cranes.find(crane => crane.Id === idx)!.HoistingTimeFilePath = await SelectFile()
    }

    const addBuilding = () => {
        if (newBuilding.Name.trim() === "") {
            alert("Building name is required")
            return
        }

        buildings.push({...newBuilding})

        // Reset form for next building
        newBuilding.Name = ""
    }

    const removeBuilding = (index: number) => {
        const buildingToRemove = buildings[index]
        if (!buildingToRemove) return

        // Store the name before removing
        const buildingName = buildingToRemove.Name

        // Remove the building
        buildings.splice(index, 1)

        // Update any cranes that were associated with this building
        cranes.forEach(crane => {
            if (crane.ForBuilding === buildingName) {
                // If there are other buildings, assign the first one
                if (buildings.length > 0) {
                    crane.ForBuilding = buildings[0].Name
                } else {
                    crane.ForBuilding = ""
                }
            }
        })
    }

    const addCrane = () => {
        if (buildings.length === 0) {
            alert("Please create at least one building first")
            currentModalStep = 'buildings'
            return
        }

        // Find the first building with a name
        const defaultBuilding = buildings.find(b => b.Name.trim() !== "")

        if (!defaultBuilding) {
            alert("Please create a building with a valid name first")
            currentModalStep = 'buildings'
            return
        }

        cranes.push({
            Id: Math.random().toString(),
            Name: "",
            HoistingTimeFilePath: "",
            ForBuilding: defaultBuilding.Name
        })
    }

    const removeCrane = (idx: string) => {
      const indexToRemove = cranes.findIndex(crane => crane.Id === idx);
      if (indexToRemove !== -1) {
        cranes.splice(indexToRemove, 1);
      }
    }

    const goToNextStep = () => {
        if (buildings.length === 0) {
            alert("Please create at least one building before proceeding")
            return
        }
        currentModalStep = 'cranes'
    }

    const goToPreviousStep = () => {
        currentModalStep = 'buildings'
    }


</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-4 grid-rows-3 ">
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
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="1" bind:value={config.AlphaHoistingPenalty}/>
    </fieldset>


    <Modal bind:isModalOpen={isOpenModal} buttonText={currentModalStep === 'buildings' ? 'Next' : 'Save'}>
      {#snippet content()}
        <div class="h-[600px] overflow-y-auto">
          {#if currentModalStep === 'buildings'}
            <!-- Step 1: Building Creation -->
            <h2 class="text-2xl font-bold mb-4">Step 1: Create Buildings</h2>
            <div class="card bg-base-100 border shadow-sm p-4 mb-6">
              <div class="grid grid-cols-3 gap-4">
                <fieldset class="fieldset flex flex-col">
                  <legend class="fieldset-legend text-base">Building Name:</legend>
                  <input type="text" class="input input-sm" placeholder="Building 1" 
                         bind:value={newBuilding.Name}/>
                </fieldset>
                <fieldset class="fieldset flex flex-col">
                  <legend class="fieldset-legend text-base">Number of Floors:</legend>
                  <input type="number" class="input input-sm" placeholder="10" 
                         bind:value={newBuilding.NumberOfFloors}/>
                </fieldset>
                <fieldset class="fieldset flex flex-col">
                  <legend class="fieldset-legend text-base">Floor Height (m):</legend>
                  <input type="number" class="input input-sm" placeholder="3.2" 
                         bind:value={newBuilding.FloorHeight}/>
                </fieldset>
              </div>
              <div class="mt-4 flex justify-end">
                <button onclick={addBuilding} class="btn btn-primary btn-sm">Add Building</button>
              </div>
            </div>

            <h3 class="text-xl font-bold mb-2">Current Buildings</h3>
            <div class="grid grid-cols-3 gap-4">
              {#if buildings.length > 0}
                {#each buildings as building, idx}
                  <div class="p-2 card bg-base-100 border shadow-sm">
                    <div class="relative card-body">
                      <div class="absolute text-lg font-bold top-0 left-1">
                        {idx + 1}.
                      </div>
                      <div class="absolute top-0 right-0 card-actions justify-end">
                        <button onclick={() => removeBuilding(idx)} aria-label="delete btn"
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

                      <div class="flex flex-col mt-4">
                        <p><strong>Name:</strong> {building.Name}</p>
                        <p><strong>Floors:</strong> {building.NumberOfFloors}</p>
                        <p><strong>Floor Height:</strong> {building.FloorHeight}m</p>
                      </div>
                    </div>
                  </div>
                {/each}
              {:else}
                <div class="col-span-3 flex items-center justify-center p-4 border rounded">
                  No buildings created yet. Please add at least one building.
                </div>
              {/if}
            </div>
          {:else}
            <!-- Step 2: Crane Creation -->
            <h2 class="text-2xl font-bold mb-4">Step 2: Create Cranes</h2>
            <div class="flex justify-between mb-4">
              <button onclick={goToPreviousStep} class="btn btn-outline btn-sm">
                Back to Buildings
              </button>
            </div>

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
                        <div class="flex">
                          <fieldset class="fieldset flex flex-col">
                            <legend class="fieldset-legend text-base">For Building:</legend>
                            <select class="select select-sm" bind:value={crane.ForBuilding}>
                              <option disabled selected></option>
                              {#each buildings as building}
                                <option value={building.Name}>{building.Name}</option>
                              {/each}
                            </select>
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
                <div class="col-span-3 flex items-center justify-center p-4 border rounded">
                    No cranes created yet. Click "Add Crane" to create one.
                </div>
              {/if}
            </div>
            <div class="mt-8 flex justify-center items-center">
              <button onclick={addCrane} class="btn btn-soft btn-primary">Add Crane</button>
            </div>
          {/if}
        </div>
      {/snippet}
      {#snippet moreButtons()}
        {#if currentModalStep === 'buildings' && buildings.length > 0}
          <button class="btn btn-primary" onclick={goToNextStep}>Next: Setup Cranes</button>
        {:else if currentModalStep === 'cranes'}
          <button class="btn btn-outline" onclick={goToPreviousStep}>Back to Buildings</button>
        {/if}
        <button class="btn btn-primary">Import Data Template</button>
      {/snippet}
    </Modal>
  </div>
  <div class="flex justify-end items-center">
    <button class="btn btn-primary" onclick={()=>isOpenModal = true}>Setup Cranes</button>
  </div>
</div>
