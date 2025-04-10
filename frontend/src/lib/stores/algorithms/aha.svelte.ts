import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IAHAConfig extends IAlgorithmSwarmConfigBase{}

export const ahaConfig = $state<IAHAConfig>({
  iterations: 300,
  population: 100,
  type: 'Swarm',
})
