<script lang="ts">
  import {type IFixedFacility, predeterminedProblemConfig} from "$lib/stores/problems";
  import Modal from "$lib/components/modal.svelte";

  const config = predeterminedProblemConfig

  let isOpenModal: boolean = $state(false);

  let fixedFacilities = $state<IFixedFacility[]>([])

  const handleOpenModal = () => {
    if (config.fixedFacilities.length === config.locationNames.length) {
      fixedFacilities = config.fixedFacilities
    } else {
      fixedFacilities = Array.from({length: config.value.numberOfLocations ?? 0}, (_, i) => ({
        LocName: `L${i + 1}`,
        FacilityName: ''
      }))
    }


    isOpenModal = true
  }

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-2 grid-rows-2 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Number of locations:</legend>
      <input type="number" class="input input-lg" placeholder="11" bind:value={config.value.numberOfLocations}/>
    </fieldset>
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Number of facilities:</legend>
      <input type="number" class="input input-lg" placeholder="11" bind:value={config.value.numberOfFacilities}/>
    </fieldset>
  </div>
  <Modal bind:isModalOpen={isOpenModal}
         buttonText="Save"
         mainActionButton={() => {
    config.setupFixedFacilities(fixedFacilities)
  }}>
    {#snippet content()}
      <div class="h-[600px] overflow-y-auto">
        <div class="grid grid-cols-6 gap-4">
          {#if config.locationNames.length > 0 && config.facilityNames.length > 0 && fixedFacilities.length === config.locationNames.length}
            {#each config.locationNames as loc, idx (loc)}
              <div class="p-2 card bg-base-100 border shadow-sm flex items-center justify-center">
                <div class="relative card-body">
                  <div class="flex flex-col">
                    <div>
                      {loc}
                    </div>
                    <div class="flex">
                      <fieldset class="fieldset flex flex-col">
                        <legend class="fieldset-legend text-base">Facility:</legend>
                        <select class="select select-sm" bind:value={fixedFacilities[idx].FacilityName}>
                          <option disabled></option>
                          {#each config.facilityNames as name (name)}
                            <option value={name}>{name}</option>
                          {/each}
                        </select>
                      </fieldset>
                    </div>
                  </div>
                </div>
              </div>
            {/each}
          {/if}
        </div>
      </div>
    {/snippet}
  </Modal>
  <div class="flex justify-end items-center">
    <button class={["btn btn-primary", {
      "btn-disabled": config.value.numberOfLocations === 0 || config.value.numberOfFacilities === 0,
    }]} onclick={handleOpenModal}>Setup located facility
    </button>
  </div>
</div>