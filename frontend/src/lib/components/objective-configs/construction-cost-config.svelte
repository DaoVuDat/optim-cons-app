<script lang="ts">
  import {SelectFile} from "$lib/wailsjs/go/main/App";
  import {constructionCostConfig} from "$lib/stores/objectives";

  const config = constructionCostConfig

</script>


<div class="p-2 w-full h-full flex flex-col justify-between">
  <div class="grid gap-2 grid-rows-2 grid-cols-2">
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg">Frequency Matrix file:</legend>

      <div class="join">
        <div>
          <label class="input input-lg validator join-item">
            <input type="text" placeholder="path://" bind:value={config.FrequencyMatrixFilePath}/>
          </label>
        </div>
        <button class="btn btn-neutral join-item btn-lg"
                onclick={async () => config.FrequencyMatrixFilePath = await SelectFile()}>Select file
        </button>
      </div>
    </fieldset>
    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg">Distance Matrix file:</legend>

      <div class="join">
        <div>
          <label class="input input-lg validator join-item">
            <input type="text" placeholder="path://" bind:value={config.DistanceMatrixFilePath}/>
          </label>
        </div>
        <button class="btn btn-neutral join-item btn-lg"
                onclick={async () => config.DistanceMatrixFilePath = await SelectFile()}>Select file
        </button>
      </div>
    </fieldset>

    <fieldset class="fieldset flex flex-col">
      <legend class="fieldset-legend text-lg ">Alpha (for Penalty):</legend>
      <input type="number" class="input input-lg" placeholder="100" bind:value={config.AlphaCCPenalty}/>
    </fieldset>
    <fieldset class="fieldset bg-base-100 border-base-300 rounded-box border p-4">
      <legend class="fieldset-legend text-lg">General QAP</legend>
      <label class="label">
        <input type="checkbox" bind:checked={config.GeneralQAP} class="toggle toggle-large checked:bg-[#422AD5] checked:text-white" />
        <span class="ml-8 text-base text-black">{config.GeneralQAP ? 'Full-matrix formulation' : 'Half-matrix formulation'}</span>
      </label>
    </fieldset>
  </div>
  <div class="flex justify-end items-center">
    <button class="btn  btn-primary">Import Data Template</button>
  </div>
</div>