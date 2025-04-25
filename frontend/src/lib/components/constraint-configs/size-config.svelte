<script lang="ts">
    import Modal from "$lib/components/modal.svelte"
    import {sizeConfig} from "$lib/stores/constraints/size.svelte";

    interface Props {
        numberOfLocations: number
        numberOfFacilities: number
    }

    const {numberOfLocations, numberOfFacilities}: Props = $props()

    const config = sizeConfig

    let isOpenModalSmallLocations = $state<boolean>(false)
    let isOpenModalLargeFacilities = $state<boolean>(false)

    const locationNames = $derived.by(() => {
        let names: string[] = []

        for (let i = 0; i < numberOfLocations; i++) {
            names.push(`L${i + 1}`)
        }

        return names
    })

    const facilityNames = $derived.by(() => {
        let names: string[] = []
        for (let i = 0; i < numberOfFacilities; i++) {
            names.push(`TF${i + 1}`)
        }
        return names
    })


    const addSmallLocation = (name: string) => {
        sizeConfig.SmallLocations.push(name)
    }

    const removeSmallLocation = (name: string) => {
        sizeConfig.SmallLocations.splice(sizeConfig.SmallLocations.indexOf(name), 1)
    }

    const addLargeFacility = (name: string) => {
        sizeConfig.LargeFacilities.push(name)
    }

    const removeLargeFacility = (name: string) => {
        sizeConfig.LargeFacilities.splice(sizeConfig.LargeFacilities.indexOf(name), 1)
    }

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-1 grid-rows-2 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Power Difference (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="10" bind:value={config.PowerDifferencePenalty}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="20000" bind:value={config.AlphaSizePenalty}/>
    </fieldset>


    <Modal bind:isModalOpen={isOpenModalSmallLocations} buttonText="Save">
      {#snippet content()}
        <div class="h-[600px] overflow-y-auto">
          <div class="grid grid-cols-3 gap-4">
            {#if sizeConfig.SmallLocations.length > 0}
              <!--  List of cranes  -->
              {#each sizeConfig.SmallLocations as smallLoc, idx (smallLoc)}
                <div class="p-2 card bg-base-100 border shadow-sm flex items-center justify-center">
                  <div class="relative card-body">
                    <div class="absolute text-lg font-bold top-0 left-1">
                      {idx + 1}.
                    </div>
                    <div class="absolute top-0 right-0 card-actions justify-end">
                      <button onclick={() => removeSmallLocation(smallLoc)} aria-label="delete btn"
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
                          <legend class="fieldset-legend text-base">Select location:</legend>
                          <select class="select select-sm" bind:value={sizeConfig.SmallLocations[idx]}>
                            {#each locationNames as name (name)}
                              <option value={name}>{name}</option>
                            {/each}
                          </select>
                        </fieldset>
                      </div>
                    </div>
                  </div>
                </div>
              {/each}
            {:else}
              <!-- No cranes -->
              <div class="col-start-2 flex items-center justify-center">
                No Small Locations
              </div>
            {/if}
          </div>
          <div class="mt-8 flex justify-center items-center">
            <button onclick={() => addSmallLocation("")} class="btn btn-soft btn-primary">Set location</button>
          </div>
        </div>
      {/snippet}
    </Modal>

    <Modal bind:isModalOpen={isOpenModalLargeFacilities} buttonText="Save">
      {#snippet content()}
        <div class="h-[600px] overflow-y-auto">
          <div class="grid grid-cols-3 gap-4">
            {#if sizeConfig.LargeFacilities.length > 0}
              <!--  List of cranes  -->
              {#each sizeConfig.LargeFacilities as fac, idx (fac)}
                <div class="p-2 card bg-base-100 border shadow-sm flex items-center justify-center">
                  <div class="relative card-body">
                    <div class="absolute text-lg font-bold top-0 left-1">
                      {idx + 1}.
                    </div>
                    <div class="absolute top-0 right-0 card-actions justify-end">
                      <button onclick={() => removeLargeFacility(fac)} aria-label="delete btn"
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
                          <legend class="fieldset-legend text-base">Select facility:</legend>
                          <select class="select select-sm" bind:value={sizeConfig.LargeFacilities[idx]}>
                            {#each facilityNames as name (name)}
                              <option value={name}>{name}</option>
                            {/each}
                          </select>
                        </fieldset>
                      </div>
                    </div>
                  </div>
                </div>
              {/each}
            {:else}
              <!-- No cranes -->
              <div class="col-start-2 flex items-center justify-center">
                No Large Facilities
              </div>
            {/if}
          </div>
          <div class="mt-8 flex justify-center items-center">
            <button onclick={() => addLargeFacility("")} class="btn btn-soft btn-primary">Set facility</button>
          </div>
        </div>
      {/snippet}
    </Modal>
  </div>
  <div class="flex justify-end items-center space-x-4">
    <button class="btn btn-primary" onclick={()=>isOpenModalSmallLocations = true}>Setup Small Locations</button>
    <button class="btn btn-primary" onclick={()=>isOpenModalLargeFacilities = true}>Setup Large Facilities</button>
  </div>
</div>