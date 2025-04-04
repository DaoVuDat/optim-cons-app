import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IMOAHAConfig extends IAlgorithmSwarmConfigBase{
  archiveSize: {
    label: 'Archive Size',
    value: number,
  };
}