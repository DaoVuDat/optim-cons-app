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

    // const updateChart = (graphsData: ResultLocationWithId[]) => {
    //
    //     const numberOfValues = graphsData[0].Value.length
    //
    //
    //     const results = graphsData.map(res => {
    //
    //
    //         const valueWithKeys = {...res.ValuesWithKey}
    //
    //         let keys = Object.keys(valueWithKeys);
    //
    //         const values = keys.map(k => {
    //             return valueWithKeys[k];
    //         });
    //
    //         // strip "Objectives"
    //         keys = keys.map(k => k.replace(/objective/gi, ""))
    //
    //         return {
    //             name: `#${parseInt(res.Id.split("-")[1]) + 1}`,
    //             valuesWithKeys: valueWithKeys,
    //             value: values,
    //             keys: keys,
    //         }
    //     })
    //
    //     const minValues = results.reduce((prev, cur) => {
    //         for (let i = 0; i < numberOfValues; i++) {
    //             if (cur.value[i] < prev[i]) {
    //                 prev[i] = cur.value[i];
    //             }
    //         }
    //
    //         return prev
    //     }, [...results[0].value])
    //
    //     const maxValues = results.reduce((prev, cur) => {
    //         for (let i = 0; i < numberOfValues; i++) {
    //             if (cur.value[i] > prev[i]) {
    //                 prev[i] = cur.value[i];
    //             }
    //         }
    //
    //         return prev
    //     }, [...results[0].value])
    //
    //     console.log(minValues, maxValues, results)
    //     let title = ``
    //     if (numberOfValues > 1) {
    //         title = `Pareto`
    //     } else {
    //         title = `Convergence`
    //     }
    //
    //     if (chartInstance) {
    //         chartInstance.clear()
    //         let options: EChartsCoreOption = {
    //             title: {
    //                 text: title,
    //                 left: 'center',
    //                 top: 'top',
    //                 textStyle: {
    //                     fontSize: 20,
    //                     fontWeight: 'bold',
    //                     color: '#000'
    //                 },
    //             },
    //             tooltip: {
    //                 trigger: 'item',
    //                 formatter: (params: any) => {
    //                     const [x, y] = params.value;
    //                     const valuesWithKey = params.data.valuesWithKeys as ValuesWithKey
    //                     const name = params.name;
    //
    //                     let tooltipString = `
    //           <div style="font-weight: bold; margin-bottom: 5px;">${name}</div>
    //         `
    //                     Object.entries(valuesWithKey).map(([k, v]) => {
    //                         tooltipString += `
    //             <div>
    //               <span style="font-weight: bold;">${k}:</span> ${v.toFixed(3)}
    //               </div>
    //             `
    //                     })
    //
    //
    //                     return tooltipString;
    //                 },
    //             },
    //
    //         }
    //
    //         if (numberOfValues > 2) {
    //             options.grid3D = {}
    //             options.xAxis3D = {
    //                 name: results[0].keys[0],
    //                 type: 'value',
    //                 min: (minValues[0] * 0.9998).toFixed(3),
    //                 max: (maxValues[0] * 1.0002).toFixed(3),
    //                 nameGap: 30
    //             }
    //             options.yAxis3D = {
    //                 name: results[0].keys[1],
    //                 type: 'value',
    //                 min: (minValues[1] * 0.9998).toFixed(3),
    //                 max: (maxValues[1] * 1.0002).toFixed(3),
    //                 nameGap: 30
    //             }
    //             options.zAxis3D = {
    //                 name: results[0].keys[2],
    //                 type: 'value',
    //                 min: (minValues[2] * 0.9998).toFixed(3),
    //                 max: (maxValues[2] * 1.0002).toFixed(3),
    //                 nameGap: 30
    //             }
    //             options.series = [
    //                 {
    //                     type: 'scatter3D',
    //                     symbolSize: 8,
    //                     data: results
    //                 }
    //             ]
    //
    //         } else if (numberOfValues == 2) {
    //             options.grid = {containLabel: true}
    //             options.xAxis = {
    //                 type: 'value',
    //                 name: results[0].keys[0],
    //                 min: (minValues[0] * 0.9998).toFixed(3),
    //                 max: (maxValues[0] * 1.0002).toFixed(3),
    //                 nameLocation: 'middle',
    //                 nameGap: 30
    //             }
    //             options.yAxis = {
    //                 type: 'value',
    //                 name: results[0].keys[1],
    //                 min: (minValues[1] * 0.9998).toFixed(3),
    //                 max: (maxValues[1] * 1.0002).toFixed(3)
    //             }
    //             options.series = [
    //                 {
    //                     data: results,
    //                     type: "scatter"
    //                 }
    //             ]
    //         } else {
    //
    //         }
    //
    //         chartInstance.setOption(options);
    //     }
    // }

    const updateChart = (graphsData: ResultLocationWithId[]) => {
        const numberOfValues = graphsData[0].Value.length

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
                convergence: res.Convergence || []
            }
        })

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
                        if (numberOfValues === 1 && params.seriesType === 'line') {
                            return `
                                <div style="font-weight: bold; margin-bottom: 5px;">Iteration ${params.dataIndex + 1}</div>
                                <div><span style="font-weight: bold;">Value:</span> ${params.value.toFixed(3)}</div>
                            `;
                        }

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

                };
                options.yAxis3D = {
                    name: results[0].keys[1],
                    type: 'value',
                    min: yBounds.min.toFixed(3),
                    max: yBounds.max.toFixed(3),
                    nameGap: 30,

                };
                options.zAxis3D = {
                    name: results[0].keys[2],
                    type: 'value',
                    min: zBounds.min.toFixed(3),
                    max: zBounds.max.toFixed(3),
                    nameGap: 30,

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

                };
                options.yAxis = {
                    type: 'value',
                    name: results[0].keys[1],
                    min: yBounds.min.toFixed(3),
                    max: yBounds.max.toFixed(3),
                    nameLocation: 'middle',
                    nameGap: 30,

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
                    // Calculate bounds for y-axis (convergence values)
                    const yBounds = calculateAxisBounds(convergenceData);

                    // Create category data for all iterations
                    const iterationCount = convergenceData.length;
                    const categoryData = Array.from({length: iterationCount}, (_, i) => i + 1);

                    options.grid = {
                        containLabel: true,
                        left: '5%',
                        right: '5%',
                        bottom: '15%'  // Increase bottom margin for dataZoom
                    };

                    // Use value type instead of category for better handling of large number of iterations
                    options.xAxis = {
                        type: 'value',
                        name: 'Iteration',
                        nameLocation: 'middle',
                        nameGap: 30,
                        min: 1,
                        max: iterationCount,
                        minInterval: 1,  // Ensure integer ticks
                        splitNumber: Math.min(10, iterationCount),  // Limit number of major ticks
                        axisLabel: {
                            formatter: (value) => {
                                // Ensure integers for iteration numbers
                                return Math.round(value).toString();
                            }
                        }
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

                    // Convert data to [x, y] pairs for value axis
                    const seriesData = convergenceData.map((value, index) => {
                        return [index + 1, value];  // [iteration number, value]
                    });

                    options.series = [
                        {
                            name: 'Convergence',
                            type: 'line',
                            data: seriesData,
                            symbolSize: 4,
                            // For large datasets, show fewer symbols
                            showSymbol: iterationCount < 100,
                            sampling: iterationCount > 100 ? 'average' : undefined,
                            connectNulls: true  // Connect across null points if any
                        }
                    ];

                    // Enhanced dataZoom for better handling of large datasets
                    options.dataZoom = [
                        {
                            type: 'inside',  // Keep the inside zoom capability
                            start: 0,
                            end: 100,
                            xAxisIndex: 0
                        }
                    ];
                }
            }

            chartInstance.setOption(options);
        }
    };

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

    const exportChart = () => {
        if (chartInstance) {
            const dataURL = chartInstance.getDataURL({
                type: 'png',
                pixelRatio: 2,
                backgroundColor: '#fff'
            });

            const link = document.createElement('a');
            link.download = 'chart-export.png';
            link.href = dataURL;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
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
        if (!values || values.length === 0) return { min: 0, max: 1, useScientific: false };

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
  <div class="w-full px-4 flex justify-between items-center">
    {#if graphsData && graphsData.length > 0}
      <button class="btn btn-primary">Export chart</button>
    {/if}
  </div>
</div>
