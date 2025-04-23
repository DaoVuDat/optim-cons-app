<script lang="ts">
  import {predeterminedProblemConfig} from "$lib/stores/problems";
  import {sizeConfig} from "$lib/stores/constraints/size.svelte";
  import Modal from "$lib/components/modal.svelte";

  const config = predeterminedProblemConfig

  interface IFixedFacility {
    LocName: string
    FacilityName : string
  }

  let fixedFacilities = $state<IFixedFacility[]>([])

  $effect(() => {
    if (locationNames.length !== fixedFacilities.length) {
      fixedFacilities = Array.from({ length: locationNames.length }, (_, i) => ({
        LocName: `L${i + 1}`,
        FacilityName: ""
      }))
    }
  })

  const locationNames = $derived.by(() => {
    let names: string[] = []

    for (let i = 0; i < config.numberOfLocations; i++) {
      names.push(`L${i + 1}`)
    }

    return names
  })

  const facilityNames = $derived.by(() => {
    let names: string[] = []
    for (let i = 0; i < config.numberOfFacilities; i++) {
      names.push(`TF${i + 1}`)
    }
    return names
  })

  let isOpenModal: boolean = $state(false);

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-2 grid-rows-2 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Number of locations:</legend>
      <input type="number" class="input input-lg" placeholder="11" bind:value={config.numberOfLocations}/>
    </fieldset>
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Number of facilities:</legend>
      <input type="number" class="input input-lg" placeholder="11" bind:value={config.numberOfFacilities}/>
    </fieldset>
  </div>
  <Modal bind:isModalOpen={isOpenModal} buttonText="Save">
    {#snippet content()}
      <div class="h-[600px] overflow-y-auto">
        <div class="grid grid-cols-3 gap-4">
          {#each locationNames as loc, idx (loc)}
            <div class="p-2 card bg-base-100 border shadow-sm flex items-center justify-center">
              <div class="relative card-body">
                <div class="flex flex-col">
                  <div>
                    {loc}
                  </div>
                  <div class="flex">
                    <fieldset class="fieldset flex flex-col">
                      <legend class="fieldset-legend text-base">Select facility:</legend>
                      <select class="select select-sm" bind:value={fixedFacilities[idx].FacilityName}>
                        <option disabled selected></option>
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
        </div>
      </div>
    {/snippet}
  </Modal>
  <div class="flex justify-end items-center">
    <button class={["btn btn-primary", {
      "btn-disabled": config.numberOfLocations > 0 && config.numberOfFacilities > 0,
    }]}>Setup located facility
    </button>
  </div>
</div>