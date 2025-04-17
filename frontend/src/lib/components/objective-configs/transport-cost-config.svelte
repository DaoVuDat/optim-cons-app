<script lang="ts">
  import {SelectFile} from "$lib/wailsjs/go/main/App";
  import {transportCostConfig} from "$lib/stores/objectives";

  const config = transportCostConfig

  const selectFile = async () => {
      config.InteractionMatrixFilePath = await SelectFile()
  }

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-rows-2 ">
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg">Facilities Interaction Matrix file:</legend>

      <div class="join">
        <div>
          <label class="input input-lg validator join-item">
            <input type="text" placeholder="path://" bind:value={config.InteractionMatrixFilePath}/>
          </label>
        </div>
        <button class="btn btn-neutral join-item btn-lg"
                onclick={() =>selectFile()}>Select file
        </button>
      </div>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="100" bind:value={config.AlphaTCPenalty}/>
    </fieldset>

  </div>
  <div class="flex justify-end items-center">
    <button class="btn  btn-primary">Import Data Template</button>
  </div>
</div>