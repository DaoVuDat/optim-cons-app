<script lang="ts">
  import {problemStore, ProblemType} from "$lib/stores/problem.svelte";
  import {SelectFile} from "$lib/wailsjs/go/main/App";
  import {ContinuousFile, continuousProblemConfig} from "$lib/stores/problems";
  import {setContext} from "svelte";

  let config = problemStore.getConfig(ProblemType.Continuous)

  const selectFile = async (field: ContinuousFile) => {

    const fileName = await SelectFile()

    switch (field) {
      case ContinuousFile.Facility:
        config.facilitiesFilePath.value = fileName;
        break
      case ContinuousFile.Phase:
        config.phasesFilePath.value = fileName;
        break;
    }
  }

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-2 grid-rows-2">
    <fieldset class="fieldset text-lg">
      <legend class="fieldset-legend">Layout length:</legend>
      <input type="text" class="input input-lg" placeholder="300" bind:value={config.length} />
    </fieldset>
    <fieldset class="fieldset">
      <legend class="fieldset-legend text-lg">Layout width:</legend>
      <input type="text" class="input input-lg" placeholder="300" bind:value={config.width}/>
    </fieldset>
    <fieldset class="fieldset">
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
    <fieldset class="fieldset">
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