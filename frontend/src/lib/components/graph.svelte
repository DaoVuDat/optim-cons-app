<script lang="ts">
  import {onMount, type Snippet} from "svelte";
  import type {ResultLocationWithId} from "../../types/result";
  import * as echarts from 'echarts';
  import type {EChartsOption, EChartsType} from "echarts";

  interface LayoutSize {
    minX: number;
    minY: number;
    maxX: number;
    maxY: number;
  }

  interface Props {
    graphData?: ResultLocationWithId
    footer?: Snippet
    layoutSize: LayoutSize
  }

  const {
    footer,
    graphData = $bindable(),
    layoutSize = $bindable()
  }: Props = $props()


  let phasesGraphData: string[][] | undefined = $derived.by(() => {
    return graphData?.Phases
  })

  let selectedPhases = $state<number>(0)

  let chartContainer: HTMLDivElement;
  let chartInstance: EChartsType;

  const updateChart = (graphData: ResultLocationWithId) => {
    // transform graphData
    const facilities = Object.values(graphData.MapLocations)
      .filter(loc => phasesGraphData![selectedPhases].includes(loc.Symbol))
      .map(loc => {
        return {
          name: loc.Symbol,
          value: [loc.Coordinate.X, loc.Coordinate.Y, loc.Length, loc.Width], // x, y, width, height
        }
      })

    const cranes = graphData.Cranes.map(loc => {
      return {
        value: [loc.Coordinate.X, loc.Coordinate.Y, loc.Length, loc.Width, loc.Radius]
      }
    })

    const maxSize = Math.max(layoutSize.maxX, layoutSize.maxY)
    const minSize = Math.min(layoutSize.minX, layoutSize.minY)

    if (chartInstance) {
      const options: EChartsOption = {
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

              console.log("input data", x, y, dataWidth, dataHeight)

              // Convert the center point to pixel coordinates
              const centerPx = api.coord([x, y]);

              // Convert data unit differences to pixel values
              // width only affects the x-axis, and height only affects the y-axis.
              const widthPx = api.size([dataWidth, 0])[0];
              const heightPx = api.size([0, dataHeight])[1];

              console.log("in pixel", centerPx, widthPx, heightPx)

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
                  style: api.style({
                    // Example style; adjust as needed
                    fill: 'rgba(0, 0, 255, 0.3)',
                    stroke: '#000',
                    lineWidth: 1
                  })
                }]
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

<div class="w-full h-full flex flex-col justify-between items-center">
  <div class="w-full mb-4 flex justify-center items-center">
    <div bind:this={chartContainer} style="width: 500px; height: 500px;"></div>
  </div>
  <div>
    Select phase
  </div>
  {#if footer}
    {@render footer()}
  {/if}
</div>

<style>

</style>