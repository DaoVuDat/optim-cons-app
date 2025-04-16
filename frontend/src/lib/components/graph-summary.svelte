<script lang="ts">
  import {onMount} from "svelte";
  import type {ResultLocationWithId, ValuesWithKey} from "../../types/result";
  import {type EChartsCoreOption, type EChartsOption, type EChartsType} from "echarts";
  import * as echarts from 'echarts';
  import _, {forEach} from 'lodash'


  interface Props {
    graphsData?: ResultLocationWithId[]
  }

  const {
    graphsData = $bindable(),
  }: Props = $props()


  $inspect(graphsData)

  let chartContainer: HTMLDivElement;
  let chartInstance: EChartsType;

  const updateChart = (graphsData: ResultLocationWithId[]) => {


    const title = `Pareto`

    const results = graphsData.map(res => {

      const valueWithKeys = {...res.ValuesWithKey}
      let keys = Object.keys(valueWithKeys);

      const values = keys.map(k => {
        return valueWithKeys[k];
      });

      // strip "Objectives"
      keys = keys.map(k => k.replace(/objective/gi, ""))

      return {
        name: `#${parseInt(res.Id.split("-")[1]) + 1}`,
        valuesWithKeys: valueWithKeys,
        value: values,
        keys: keys,
      }
    })

    const numberOfValues = results[0].value.length
    const minValues = results.reduce((prev, cur) => {
      const firstVal = cur.value[0]
      const secondVal = cur.value[1]
      if (firstVal < prev[0]) {
        prev[0] = firstVal;
      }

      if (secondVal < prev[1]) {
        prev[1] = secondVal;
      }

      return prev
    }, [...results[0].value])

    const maxValues = results.reduce((prev, cur) => {
      const firstVal = cur.value[0]
      const secondVal = cur.value[1]
      if (firstVal > prev[0]) {
        prev[0] = firstVal;
      }

      if (secondVal > prev[1]) {
        prev[1] = secondVal;
      }

      return prev
    }, [...results[0].value])

    console.log(minValues, maxValues, results)

    if (chartInstance) {
      chartInstance.clear()
      let options: EChartsCoreOption = {
        title: {
          text: title,
          left: 'center',
          top: 'top',
          textStyle: {
            fontSize: 20,
            fontWeight: 'bold',
            color: '#000'
          },
        },
        tooltip: {
          trigger: 'item',
          formatter: (params: any) => {
            const [x, y] = params.value;
            const valuesWithKey = params.data.valuesWithKeys as ValuesWithKey
            const name = params.name;

            let tooltipString = `
              <div style="font-weight: bold; margin-bottom: 5px;">${name}</div>
            `
            Object.entries(valuesWithKey).map(([k, v]) => {
              tooltipString += `
                <div>
                  <span style="font-weight: bold;">${k}:</span> ${v.toFixed(3)}
                  </div>
                `
            })


            return tooltipString;
          },
        },
        grid: {containLabel: true},
      }

      if (numberOfValues > 2) {
        options = {}
      } else {
        options.xAxis = {
          type: 'value',
          name: results[0].keys[0],
          min: (minValues[0] * 0.9998).toFixed(3),
          max: (maxValues[0] * 1.0002).toFixed(3),
          nameLocation: 'middle',
          nameGap: 30
        }
        options.yAxis = {
          type: 'value',
          name: results[0].keys[1],
          min: (minValues[1] * 0.9998).toFixed(3),
          max: (maxValues[1] * 1.0002).toFixed(3)
        }
        options.series = [
          {
            data: results,
            type: "scatter"
          }
        ]
      }
      // Pareto


      chartInstance.setOption(options);
    }
  }

  $effect(() => {
    if (graphsData && graphsData.length > 0) {
      updateChart(graphsData)
    }
  })

  onMount(() => {
    chartInstance = echarts.init(chartContainer);
    if (graphsData && graphsData.length > 0) {
      updateChart(graphsData)
    }
  });

</script>

<div class="w-full h-full flex flex-col justify-center items-center">
  <div class="w-full flex justify-center items-center">
    <div bind:this={chartContainer} class="w-full" style="height: 510px;"></div>
  </div>
  <div class="w-full px-4 flex justify-between items-center">
    {#if graphsData && graphsData.length > 0}
      <button class="btn btn-primary">Export chart</button>
    {/if}
  </div>
</div>
