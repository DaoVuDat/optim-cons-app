<script lang="ts">
  import {SelectFile} from "$lib/wailsjs/go/main/App";
  import {riskConfig} from "$lib/stores/objectives";

  const config = riskConfig

  const selectFile = async () => {
    config.HazardInteractionMatrixFilePath = await SelectFile()
  }

</script>

<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-cols-2 ">
    <fieldset class="fieldset w-full flex flex-col">
      <legend class="fieldset-legend text-lg">Delta:</legend>
      <input type="number" class="input input-lg" placeholder="10" bind:value={config.Delta}/>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg">Hazard Interaction Matrix file:</legend>

      <div class="join">
        <div>
          <label class="input input-lg validator join-item">
            <input type="text" placeholder="path://" bind:value={config.HazardInteractionMatrixFilePath}/>
          </label>
        </div>
        <button class="btn btn-neutral join-item btn-lg"
                onclick={() =>selectFile()}>Select file
        </button>
      </div>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="18.75" bind:value={config.AlphaRiskPenalty}/>
    </fieldset>

  </div>
  <div class="flex justify-end items-center">
    <button class="btn  btn-primary">Import Data Template</button>
  </div>
</div>