import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IAHAConfig extends IAlgorithmSwarmConfigBase{}

export const ahaConfig = $state<IAHAConfig>({
  iterations: {
    label: 'Iterations',
    value: 300,
  },
  population: {
    label: 'Population',
    value: 100,
  },
  type: 'Swarm',
})
