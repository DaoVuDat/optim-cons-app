import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IMOGWOConfig extends IAlgorithmSwarmConfigBase {
  aParam: number,
  archiveSize: number,
  numberOfGrids: number,
  alpha: number,
  beta: number,
  gamma: number,
}

export const mogwoConfig = $state<IMOGWOConfig>({
  iterations: 300,
  population: 100,
  aParam: 2,
  archiveSize: 100,
  numberOfGrids: 10,
  alpha: 0.1,
  beta: 4,
  gamma: 2,
  type: 'Swarm',
})
