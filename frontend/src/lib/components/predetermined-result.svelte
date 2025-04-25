<script lang="ts">
  import type {ResultLocationWithId} from "../../types/result";
  import clsx from "clsx";
  import {roundNDecimal} from "$lib/utils/rounding";


  interface Props {
    graphData?: ResultLocationWithId
  }

  const {
    graphData = $bindable(),
  }: Props = $props()

  $inspect(graphData)
</script>

<div class="w-full h-full flex flex-col items-center py-6">
  {#if graphData}
    <div class="pl-6 w-full">
      <div class="grid grid-rows-2 grid-cols-2 gap-1">
        {#each Object.entries(graphData.ValuesWithKey) as [k, v]}
          <p>
              <span class="font-bold text-base text-gray-700">
                {k.replace(/objective/gi, "")}:
              </span>
            <span class="text-sm text-gray-600">
                {roundNDecimal(v, 3)}
              </span>
          </p>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Scrollable section -->
  <div class="mt-8 w-full max-w-6xl h-full overflow-y-auto space-y-8 px-4">
    <!-- Location cards -->
    {#if graphData}
      <div class="grid grid-cols-3 gap-6">
        {#each Object.values(graphData.MapLocations) as loc (loc.Symbol)}
          <div class="bg-white rounded-2xl shadow-md p-6 hover:shadow-lg transition duration-300">
            <h3 class="text-lg font-semibold text-gray-800">{loc.Symbol} assigned to <span class="font-medium text-blue-600">{loc.IsLocatedAt}</span></h3>
          </div>
        {/each}

      </div>
    {/if}



  </div>
</div>