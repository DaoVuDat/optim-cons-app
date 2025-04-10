import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IMOAHAConfig extends IAlgorithmSwarmConfigBase{
  archiveSize: number,
}

export const moahaConfig = $state<IMOAHAConfig>({
  iterations: 300,
  population: 100,
  archiveSize: 100,
  type: 'Swarm',
})