import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IOMOAHAConfig extends IAlgorithmSwarmConfigBase{
  archiveSize: number,
}

export const omoahaConfig = $state<IOMOAHAConfig>({
  iterations: 300,
  population: 100,
  archiveSize: 100,
  type: 'Swarm',
})