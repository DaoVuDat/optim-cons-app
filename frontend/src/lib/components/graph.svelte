<script lang="ts">
  import {onMount} from "svelte";
  import type {ResultLocationWithId} from "../../types/result";
  import type {EChartsOption, EChartsType} from "echarts";
  import * as echarts from 'echarts';
  import {generateLabelFriendlyColors} from "$lib/utils/generateColors";
  import clsx from "clsx";
  import type {CustomElementOption} from "echarts/types/src/chart/custom/CustomSeries";
  import {roundNDecimal} from "$lib/utils/rounding";

  interface LayoutSize {
    minX: number;
    minY: number;
    maxX: number;
    maxY: number;
  }

  interface Props {
    graphData?: ResultLocationWithId
    layoutSize: LayoutSize
    gridSize?: number,
    useGrid?: boolean,
  }

  const {
    graphData = $bindable(),
    layoutSize = $bindable(),
    useGrid = true,
    gridSize = 1,
  }: Props = $props()


  const numberOfGridCols = $derived.by(() => {
    if (!graphData) return
    return Math.ceil(graphData.Value.length / 2)
  })

  const phasesGraphData: string[][] | undefined = $derived.by(() => {
    if (!graphData) return
    // check is whether static or not

    const phases = [...graphData.Phases]

    if (phases.length > 1) {
      // add all facilities
      const allFacilities: string[] = Object.values(graphData.MapLocations).map(loc => loc.Symbol)
      phases.push(allFacilities)
    }

    return phases
  })

  let selectedPhases = $state<number>(0)

  const generatedColors = $derived.by(() => {
    if (!graphData) return

    return generateLabelFriendlyColors(Object.values(graphData.MapLocations).length, true)
  })

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
            name: loc.Name,
          }
        }
      })

    const maxSize = Math.max(layoutSize.maxX, layoutSize.maxY)
    const minSize = Math.min(layoutSize.minX, layoutSize.minY)

    const title = `Construction Layout #${parseInt(graphData.Id.split("-")[1]) + 1}`

    //
    let axisXAndYOptionsForGrid = {}
    if (useGrid) {
      const numberOfGrid = gridSize
      const fixedMajorGrid = 6
      axisXAndYOptionsForGrid = {
        splitLine: { show: true },
        // for example, auto‐choose ~6 major segments:
        // splitNumber: fixedMajorGrid,
        // // minor ticks between each major segment:
        // minorTick: {
        //   show: true,
        //   // number of minor segments between majors:
        //   // e.g. (majorSpan / desiredUnit) = (120/6) / 1 = 20
        //   splitNumber: (maxSize / fixedMajorGrid) / numberOfGrid
        // },
        // draw the minor (sub‑grid) lines
        minorSplitLine: {
          show: true,
          lineStyle: {
            type: 'dashed',
            width: 1,
            color: '#ddd'
          }
        }
      }
    }

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
          ...axisXAndYOptionsForGrid
        },
        yAxis: {
          max: maxSize,
          min: minSize,
          type: 'value',
          ...axisXAndYOptionsForGrid
        },
        tooltip: {
          trigger: 'item',
          formatter: function (params) {
            const facilityInfo = params.data.facilityInfo;
            return `
                          <div style="font-weight: bold; margin-bottom: 5px;">${params.name} - ${facilityInfo.name}</div>
                          <div><span style="font-weight: bold;">Position:</span> (${params.value[0].toFixed(2)}, ${params.value[1].toFixed(2)})</div>
                          <div><span style="font-weight: bold;">Dimensions:</span> ${params.value[2]} × ${params.value[3]}</div>
                          <div><span style="font-weight: bold;">Rotated:</span> ${facilityInfo.rotation}</div>
                          <div><span style="font-weight: bold;">Fixed:</span> ${facilityInfo.isFixed}</div>
                        `;
          }
        },
        // Enable zooming
        toolbox: {
          feature: {
            dataZoom: {show: true},
            restore: {show: true},
          }
        },
        // Enable mouse wheel zooming
        dataZoom: [
          {
            id: 'dataZoomX',
            type: 'inside',
            xAxisIndex: [0],  // Apply to first X axis
            filterMode: 'filter',
            start: 0,         // Initial range start at 0%
            end: 100          // Initial range end at 100%
          },
          {
            id: 'dataZoomY',
            type: 'inside',
            yAxisIndex: [0],  // Apply to first Y axis
            filterMode: 'filter',
            start: 0,         // Initial range start at 0%
            end: 100          // Initial range end at 100%
          }
        ],
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

              const color = generatedColors![facilityIdx]


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


              let children: CustomElementOption[] = [{
                type: 'rect',
                z: 2,
                shape: {
                  // Offset so the rectangle is centered (top-left of rect = -half width, -half height)
                  x: -widthPx / 2,
                  y: -heightPx / 2,
                  width: widthPx,
                  height: heightPx
                },
                style: style,
                enterFrom: {
                  style: {opacity: 0},
                },
                leaveTo: {
                  style: {opacity: 0},
                },
                enterAnimation: {
                  duration: 300
                },
                leaveAnimation: {
                  duration: 300
                }
              }, {
                type: 'text',
                z: 3,
                style: {
                  text: name,
                  textAlign: 'center',
                  textVerticalAlign: 'middle',
                  fill: '#000',
                  font: '13px sans-serif'
                },
                enterFrom: {
                  style: {opacity: 0},
                },
                leaveTo: {
                  style: {opacity: 0},
                },
                enterAnimation: {
                  duration: 200
                },
                leaveAnimation: {
                  duration: 200
                }
              }]

              // is crane -> draw circle
              if (graphData.Cranes) {
                const crane = graphData.Cranes.find(c => c.CraneSymbol === name)
                let circle: CustomElementOption | undefined
                if (crane) {
                  const radiusPx = api.size([crane.Radius, 0])[0];

                  circle = {
                    type: 'circle',
                    z: 1,
                    shape: {
                      // Offset so the rectangle is centered (top-left of rect = -half width, -half height)
                      x: -widthPx / 2,
                      y: -heightPx / 2,
                      r: radiusPx,
                    },
                    style: {
                      fill: 'none',
                      stroke: color,
                      lineWidth: 2
                    },
                    enterFrom: {
                      style: {opacity: 0},
                    },
                    leaveTo: {
                      style: {opacity: 0},
                    },
                    enterAnimation: {
                      duration: 200
                    },
                    leaveAnimation: {
                      duration: 200
                    }
                  }

                  children.push(circle)
                }
              }


              // Return a group that contains a rectangle.
              // We position the rectangle so that its center is at the group's origin.
              return {
                type: 'group',
                position: centerPx, // positions the group at the center point (in pixels)
                children,

              }
            },
          },
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
  <div class="w-full px-4 flex justify-between items-center">
    {#if graphData}
      <div class="grid grid-rows-2 grid-cols-{numberOfGridCols} gap-1 flex-2/3">
        {#each Object.entries(graphData.ValuesWithKey) as [k, v]}
          <p> <span class="font-bold text-base">
            {k.replace(/objective/gi, "")}:
          </span> <span class="text-sm">{roundNDecimal(v, 3)}</span></p>
        {/each}
      </div>
    {/if}
    {#if phasesGraphData && phasesGraphData.length > 1}
      <div class="join flex-1/3 justify-end">
        <button class={clsx("join-item btn", {
        "btn-disabled": selectedPhases === 0,
    })} onclick={() => selectedPhases--}>«
        </button>
        <button class="join-item btn w-48">
          {selectedPhases === phasesGraphData.length - 1 ? "All" : `Phase / Time Interval: ${selectedPhases + 1}`}
        </button>
        <button class={clsx("join-item btn", {
        "btn-disabled":  !graphData ||  selectedPhases === phasesGraphData.length - 1,
    })} onclick={() => selectedPhases++}>»
        </button>
      </div>
    {/if}
  </div>
</div>
