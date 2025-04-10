import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IGWOConfig extends IAlgorithmSwarmConfigBase {
  aParam: number
}


export const gwoConfig = $state<IGWOConfig>({
  iterations: 300,
  population: 100,
  aParam: 2,
  type: 'Swarm',
})

