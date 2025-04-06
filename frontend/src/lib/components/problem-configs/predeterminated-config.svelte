<script lang="ts">
  import {problemStore} from "$lib/stores/problem.svelte";
  import {SelectFile} from "$lib/wailsjs/go/main/App";
  import {objectives} from "$lib/wailsjs/go/models.js";
  import {PredeteriminatedFile} from "$lib/stores/problems";

  const config = problemStore.getConfig(objectives.ProblemType.PredeterminedConstructionLayout)

  const selectFile = async (field: PredeteriminatedFile) => {

    const fileName = await SelectFile()

    switch (field) {
      case PredeteriminatedFile.Predeterminated:
        config.predeterminatedLocationsFilePath.value = fileName;
        break
      default:
        break
    }
  }

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-2 grid-rows-2 ">
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg">Predeterminated Locations file:</legend>
      <div class="join">
        <div>
          <label class="input validator join-item">
            <input type="text" placeholder="path://" bind:value={config.predeterminatedLocationsFilePath.value} />
          </label>
        </div>
        <button class="btn btn-neutral join-item" onclick={() =>selectFile(config.predeterminatedLocationsFilePath.label)}>Select file</button>
      </div>
    </fieldset>
  </div>
  <div class="flex justify-end items-center">
    <button class="btn btn-primary">Import Data Template</button>
  </div>
</div>