<script lang="ts">
    import {onMount, type Snippet} from "svelte";
    import type {ResultLocationWithId} from "../../types/result";
    import * as echarts from 'echarts';
    import type {EChartsOption, EChartsType} from "echarts";
    import {derived} from "svelte/store";
    import {generateLabelFriendlyColors} from "$lib/utils/generateColors";
    import clsx from "clsx";

    interface LayoutSize {
        minX: number;
        minY: number;
        maxX: number;
        maxY: number;
    }

    interface Props {
        graphData?: ResultLocationWithId
        layoutSize: LayoutSize
    }

    const {
        graphData = $bindable(),
        layoutSize = $bindable()
    }: Props = $props()


    let phasesGraphData: string[][] | undefined = $derived.by(() => {
        return graphData?.Phases
    })

    let selectedPhases = $state<number>(0)

    let generatedColors = $derived.by(() => {
            const color = generateLabelFriendlyColors(Object.values(graphData?.MapLocations).length, true)

            return color
        }
    )

    $inspect(graphData)

    let chartContainer: HTMLDivElement;
    let chartInstance: EChartsType;

    const updateChart = (graphData: ResultLocationWithId) => {
        // transform graphData
        const facilities = Object.values(graphData.MapLocations)
            .map((loc, idx) => ({
                ...loc,
                idx: idx
            }))
            .filter(loc => phasesGraphData![selectedPhases].includes(loc.Symbol))
            .map(loc => {
                return {
                    idx: loc.idx,
                    name: loc.Symbol,
                    value: [loc.Coordinate.X, loc.Coordinate.Y, loc.Length, loc.Width], // x, y, width, height
                    facilityInfo: {
                        rotation: loc.Rotation ? 'Yes' : 'No',
                        isFixed: loc.IsFixed ? 'Yes' : 'No',
                    }
                }
            })


        const cranes = graphData.Cranes.map(loc => {
            return {
                value: [loc.Coordinate.X, loc.Coordinate.Y, loc.Length, loc.Width, loc.Radius]
            }
        })

        const craneSymbols = graphData.Cranes.map(crane => crane.CraneSymbol)

        const maxSize = Math.max(layoutSize.maxX, layoutSize.maxY) + 2
        const minSize = Math.min(layoutSize.minX, layoutSize.minY)

        const title = `Construction Layout #${parseInt(graphData.Id.split("-")[1]) + 1}`

        if (chartInstance) {
            const options: EChartsOption = {
                title: {
                    text: title,
                    left: 'center',
                    top: 'top',
                    textStyle: {
                        fontSize: 20,
                        fontWeight: 'bold',
                        color: '#000'
                    }
                },
                xAxis: {
                    max: maxSize,
                    min: minSize,
                    type: 'value',

                },
                yAxis: {
                    max: maxSize,
                    min: minSize,
                    type: 'value'
                },
                tooltip: {
                    trigger: 'item',
                    formatter: function (params) {
                        const facilityInfo = params.data.facilityInfo;
                        return `
                          <div style="font-weight: bold; margin-bottom: 5px;">${params.name}</div>
                          <div>Position: (${params.value[0].toFixed(2)}, ${params.value[1].toFixed(2)})</div>
                          <div>Dimensions: ${params.value[2]} × ${params.value[3]}</div>
                          <div>Rotated: ${facilityInfo.rotation}</div>
                          <div>Fixed: ${facilityInfo.isFixed}</div>
                        `;
                    }
                },
                series: [
                    {
                        // Custom series for rectangles (facilities)
                        data: facilities,
                        type: 'custom',
                        renderItem: function (params, api) {
                            // Extract values based on field order:
                            // [x, y, width, height]
                            const x = api.value(0);
                            const y = api.value(1);
                            const dataWidth = api.value(2);
                            const dataHeight = api.value(3);

                            const idx = params.dataIndex
                            const name = facilities[idx].name
                            const facilityIdx = facilities[idx].idx

                            const color = generatedColors[facilityIdx]

                            let style = {
                                fill: color,
                                stroke: '#000',
                                lineWidth: 1,
                            };

                            // Convert the center point to pixel coordinates
                            const centerPx = api.coord([x, y]);

                            // Convert data unit differences to pixel values
                            // width only affects the x-axis, and height only affects the y-axis.
                            const widthPx = api.size([dataWidth, 0])[0];
                            const heightPx = api.size([0, dataHeight])[1];

                            // Return a group that contains a rectangle.
                            // We position the rectangle so that its center is at the group's origin.
                            return {
                                type: 'group',
                                position: centerPx, // positions the group at the center point (in pixels)
                                children: [{
                                    type: 'rect',
                                    shape: {
                                        // Offset so the rectangle is centered (top-left of rect = -half width, -half height)
                                        x: -widthPx / 2,
                                        y: -heightPx / 2,
                                        width: widthPx,
                                        height: heightPx
                                    },

                                    style: style
                                },
                                    {
                                        type: 'text',
                                        style: {
                                            text: name,
                                            textAlign: 'center',
                                            textVerticalAlign: 'middle',
                                            fill: '#000',
                                            font: '13px sans-serif'
                                        }
                                    }
                                ]
                            }
                        },
                    },
                    // {
                    //   // draw circle
                    //   datasetIndex: cranes,
                    //   type: 'custom',
                    //
                    // }
                ]
            };

            chartInstance.setOption(options);
        }
    }

    $effect(() => {
        if (graphData) {

            updateChart(graphData)
        }
    })

    onMount(() => {
        chartInstance = echarts.init(chartContainer);
        if (graphData) {
            updateChart(graphData)
        }

    });

</script>

<div class="w-full h-full flex flex-col justify-center items-center">
  <div class="w-full flex justify-center items-center">
    <div bind:this={chartContainer} style="width: 500px; height: 500px;"></div>
  </div>
  <div class="join">
    <button class={clsx("join-item btn", {
        "btn-disabled": selectedPhases === 0,
    })} onclick={() => selectedPhases--}>«
    </button>
    <button class="join-item btn">Phase {selectedPhases + 1}</button>
    <button class={clsx("join-item btn", {
        "btn-disabled":  !graphData ||  selectedPhases === graphData.Phases.length - 1,
    })} onclick={() => selectedPhases++}>»
    </button>
  </div>
</div>

<style>

</style>