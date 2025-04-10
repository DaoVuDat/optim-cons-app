<script lang="ts">
  import {SelectFile} from "$lib/wailsjs/go/main/App";
  import {GridFile, gridProblemConfig} from "$lib/stores/problems";

  const config = gridProblemConfig

  const selectFile = async (field: GridFile) => {

    const fileName = await SelectFile()

    switch (field) {
      case GridFile.Facility:
        config.facilitiesFilePath.value = fileName;
        break
      case GridFile.Phase:
        config.phasesFilePath.value = fileName;
        break;
    }
  }

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-2 grid-rows-3 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Layout length:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.length} />
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Layout width:</legend>
      <input type="number" class="input input-lg" placeholder="300" bind:value={config.width}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Grid size:</legend>
      <input type="text" class="input input-lg" placeholder="2x2" bind:value={config.gridSize}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col row-start-3">
      <legend class="fieldset-legend text-lg">Facilities file:</legend>
      <div class="join">
        <div>
          <label class="input validator join-item">
            <input type="text" placeholder="path://" bind:value={config.facilitiesFilePath.value}/>
          </label>
        </div>
        <button class="btn btn-neutral join-item" onclick={() =>selectFile(config.facilitiesFilePath.label)}>Select file</button>
      </div>
    </fieldset>
    <fieldset class="fieldset flex flex-col row-start-3">
      <legend class="fieldset-legend text-lg">Static / Phase / Dynamic file:</legend>
      <div class="join">
        <div>
          <label class="input validator join-item">
            <input type="text" placeholder="path://" bind:value={config.phasesFilePath.value} />
          </label>
        </div>
        <button class="btn btn-neutral join-item" onclick={() =>selectFile(config.phasesFilePath.label)}>Select file</button>
      </div>
    </fieldset>
  </div>
  <div class="flex justify-end items-center">
    <button class="btn btn-primary">Import Data Template</button>
  </div>
</div>