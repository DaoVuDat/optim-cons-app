import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IMOAHAConfig extends IAlgorithmSwarmConfigBase{
  archiveSize: {
    label: 'Archive Size',
    value: number,
  };
}

export const moahaConfig = $state<IMOAHAConfig>({
  iterations: {
    label: 'Iterations',
    value: 300,
  },
  population: {
    label: 'Population',
    value: 100,
  },
  archiveSize: {
    label: "Archive Size",
    value: 100,
  },
  type: 'Swarm',
})