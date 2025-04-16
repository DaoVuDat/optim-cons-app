<script lang="ts">
    import {onMount} from "svelte";
    import type {ResultLocationWithId, ValuesWithKey} from "../../types/result";
    import {type EChartsCoreOption, type EChartsType} from "echarts";
    import 'echarts-gl';
    import * as echarts from 'echarts';


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

            // testing 3rd objective
            const randomVal = Math.random()

            const valueWithKeys = {...res.ValuesWithKey}

            // testing 3rd objective
            valueWithKeys.random = randomVal

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

            // testing 3rd objective
            prev[2] = 0

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

            // testing 3rd objective
            prev[2] = 1

            return prev
        }, [...results[0].value])

        console.log(minValues, maxValues, results)

        function generateData() {
            let data = [];
            let dataCount = 100;

            for (let i = 0; i < dataCount; i++) {
                let x = Math.random() * 10 - 5;
                let y = Math.random() * 10 - 5;
                let z = Math.random() * 10 - 5;

                // Calculate a value based on distance to origin
                let value = Math.sqrt(x * x + y * y + z * z);

                data.push([x, y, z, value]);
            }
            return data;
        }

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

            }

            if (numberOfValues > 2) {
                options.grid3D = {}
                options.xAxis3D = {
                    name: results[0].keys[0],
                    type: 'value',
                    min: (minValues[0] * 0.9998).toFixed(3),
                    max: (maxValues[0] * 1.0002).toFixed(3),
                }
                options.yAxis3D = {
                    name: results[0].keys[1],
                    type: 'value',
                    min: (minValues[1] * 0.9998).toFixed(3),
                    max: (maxValues[1] * 1.0002).toFixed(3),
                }
                options.zAxis3D = {
                    name: results[0].keys[2],
                    type: 'value',
                    min: 0,
                    max: 1,
                }
                options.series = [
                    {
                        type: 'scatter3D',
                        symbolSize: 8,
                        data: results
                    }
                ]

            } else {
                options.grid = {containLabel: true}
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
