<script lang="ts">
  import {onMount} from "svelte";
  import type {ResultLocationWithId, ValuesWithKey} from "../../types/result";
  import {type EChartsCoreOption, type EChartsType} from "echarts";
  import 'echarts-gl';
  import * as echarts from 'echarts';
  import {toast} from "@zerodevx/svelte-toast";
  import {errorOpts, successOpts} from "$lib/utils/toast-opts";
  import {SaveChartImage} from "$lib/wailsjs/go/main/App";


  interface Props {
    graphsData?: ResultLocationWithId[]
    convergence?: number[]
  }

  const {
    graphsData = $bindable(),
    convergence,
  }: Props = $props()


  // $inspect(graphsData)

  let numberOfValues = $derived(graphsData![0].Value.length)
  let modeAxis = $state<boolean>(false);

  let chartContainer: HTMLDivElement;
  let chartInstance: EChartsType;

  const updateChart = (graphsData: ResultLocationWithId[], mode: boolean) => {
    // const numberOfValues = graphsData[0].Value.length

    const results = graphsData.map(res => {
      const valueWithKeys = {...res.ValuesWithKey}
      let keys = Object.keys(valueWithKeys).sort((a, b) => mode ? a.localeCompare(b) : b.localeCompare(a));
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
        convergence: res.Convergence || []
      }
    })
    console.log(results)

    let title = '';
    if (numberOfValues > 1) {
      title = 'Pareto';
    } else {
      title = 'Convergence';
    }

    if (chartInstance) {
      chartInstance.clear();
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
            const valuesWithKey = params.data.valuesWithKeys as ValuesWithKey;
            const name = params.name;

            let tooltipString = `
                            <div style="font-weight: bold; margin-bottom: 5px;">${name}</div>
                        `;

            Object.entries(valuesWithKey).map(([k, v]) => {
              tooltipString += `
                                <div>
                                    <span style="font-weight: bold;">${k}:</span> ${v.toFixed(3)}
                                </div>
                            `;
            });

            return tooltipString;
          },
        },
      };

      if (numberOfValues > 2) {
        // 3D Pareto
        // Extract values for each dimension
        const xValues = results.map(r => r.value[0]);
        const yValues = results.map(r => r.value[1]);
        const zValues = results.map(r => r.value[2]);

        // Calculate bounds for each axis
        const xBounds = calculateAxisBounds(xValues);
        const yBounds = calculateAxisBounds(yValues);
        const zBounds = calculateAxisBounds(zValues);

        options.grid3D = {};
        options.xAxis3D = {
          name: results[0].keys[0],
          type: 'value',
          min: xBounds.min.toFixed(3),
          max: xBounds.max.toFixed(3),
          nameGap: 30,
          axisLabel: {
            formatter: formatAxisLabel
          }
        };
        options.yAxis3D = {
          name: results[0].keys[1],
          type: 'value',
          min: yBounds.min.toFixed(3),
          max: yBounds.max.toFixed(3),
          nameGap: 30,
          axisLabel: {
            formatter: formatAxisLabel
          }
        };
        options.zAxis3D = {
          name: results[0].keys[2],
          type: 'value',
          min: zBounds.min.toFixed(3),
          max: zBounds.max.toFixed(3),
          nameGap: 30,
          axisLabel: {
            formatter: formatAxisLabel
          }
        };
        options.series = [
          {
            type: 'scatter3D',
            symbolSize: 8,
            data: results
          }
        ];
      } else if (numberOfValues == 2) {
        // 2D Pareto
        // Extract values for each dimension
        const xValues = results.map(r => r.value[0]);
        const yValues = results.map(r => r.value[1]);

        // Calculate bounds for each axis
        const xBounds = calculateAxisBounds(xValues);
        const yBounds = calculateAxisBounds(yValues);

        options.grid = {containLabel: true};
        options.xAxis = {
          type: 'value',
          name: results[0].keys[0],
          min: xBounds.min.toFixed(3),
          max: xBounds.max.toFixed(3),
          nameLocation: 'middle',
          nameGap: 30,
          axisLabel: {
            formatter: formatAxisLabel
          }
        };
        options.yAxis = {
          type: 'value',
          name: results[0].keys[1],
          min: yBounds.min.toFixed(3),
          max: yBounds.max.toFixed(3),
          nameLocation: 'middle',
          nameGap: 70,
          axisLabel: {
            formatter: formatAxisLabel
          }
        };
        options.series = [
          {
            data: results,
            type: "scatter"
          }
        ];
      } else {
        // Convergence chart (when numberOfValues == 1)
        const convergenceData = results[0].convergence;
        if (convergenceData && convergenceData.length > 0) {
          options.tooltip = {
            trigger: 'axis',
            formatter: (params: any) => {
              const iterations = params[0].data[0]
              const value = params[0].data[1]

              return `<div style="font-weight: bold;">Iteration ${iterations}
                      </div>
                      <div>${value.toFixed(4)}</div>
                    `
            },
          }

          // Calculate bounds for y-axis (convergence values)
          const yBounds = calculateAxisBounds(convergenceData);

          // Create category data for all iterations
          const iterationCount = convergenceData.length;

          // Use value type instead of category for better handling of large number of iterations
          options.xAxis = {
            type: 'value',
            name: 'Iteration',
            nameLocation: 'middle',
            nameGap: 30,
            min: 1,
            max: iterationCount,
          };

          options.yAxis = {
            type: 'value',
            name: results[0].keys[0],
            nameLocation: 'middle',
            nameGap: 60,
            min: yBounds.min.toFixed(3),
            max: yBounds.max.toFixed(3),
            axisLabel: {
              formatter: formatAxisLabel
            }
          };

          const seriesData = convergenceData.map((value, index) => {
            return [index + 1, value];  // [iteration number, value]
          });
          options.series = [
            {
              name: 'Convergence',
              type: 'line',
              data: seriesData,
              smooth: true,
              showSymbol: false,
              lineStyle: {
                width: 2
              }
            }
          ];
        }
      }

      chartInstance.setOption(options);
    }
  };

  const handleChangeAxis = () => {
    modeAxis = !modeAxis
  }

  $effect(() => {
    if (graphsData && graphsData.length > 0) {
      updateChart(graphsData, modeAxis)
    }
  })

  onMount(() => {
    chartInstance = echarts.init(chartContainer);
    if (graphsData && graphsData.length > 0) {
      updateChart(graphsData, modeAxis)
    }
  });

  const exportChart = async () => {
    if (chartInstance) {
      const dataURL = chartInstance.getDataURL({
        type: 'png',
        pixelRatio: 2,
        backgroundColor: '#fff'
      });
      try {
        // Call the Wails backend to save the chart
        const savedPath = await SaveChartImage(dataURL);

        if (savedPath) {
          // Show success notification using Svelte toast
          toast.push(`Chart has been saved to: ${savedPath}`, {
            theme: successOpts
          });
        }
      } catch (error) {
        console.error("Error saving chart:", error);
        toast.push(`Error saving chart: ${error as string}`, {
          theme: errorOpts
        });
        // Show error notification using Svelte toast
      }
    }
  };

  // Custom formatter function to handle zero values properly
  const formatAxisLabel = (value: number) => {
    if (value === 0) return '0';
    if (Math.abs(value) < 0.001 || Math.abs(value) >= 10000) {
      // Use scientific notation for very small or very large numbers
      return value.toExponential(2);
    }
    // Use fixed decimal notation for normal ranges
    return value.toFixed(2);
  };

  // Helper function to determine if scientific notation should be used
  const shouldUseScientific = (min: number, max: number): boolean => {
    // Use scientific notation if values are very large or very small
    const absMin = Math.abs(min);
    const absMax = Math.abs(max);
    return (absMax >= 10000 || absMin >= 10000 || (absMax > 0 && absMax < 0.01) || (absMin > 0 && absMin < 0.01));
  };

  const calculateAxisBounds = (values: number[]) => {
    if (!values || values.length === 0) return {min: 0, max: 1, useScientific: false};

    const min = Math.min(...values);
    const max = Math.max(...values);
    const range = max - min;

    // Add 2% padding on each side
    const minBound = min - (range * 0.02);
    const maxBound = max + (range * 0.02);

    // Determine if we should use scientific notation
    const useScientific = shouldUseScientific(minBound, maxBound);

    return {
      min: minBound,
      max: maxBound,
      useScientific
    };
  };

</script>

<div class="w-full h-full flex flex-col justify-center items-center">
  <div class="w-full flex justify-center items-center">
    <div bind:this={chartContainer} style="height: 510px; width: 100%;"></div>
  </div>
  <div class="w-full px-4 flex justify-end items-center space-x-4">
    {#if graphsData && graphsData.length > 0}
      {#if numberOfValues === 2}
        <button class="btn btn-primary" onclick={handleChangeAxis}>Change Axis</button>
      {/if}
      <button class="btn btn-primary" onclick={exportChart}>Export chart</button>
    {/if}
  </div>
</div>
