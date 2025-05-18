import type {IAlgorithmSwarmConfigBase} from "$lib/stores/algorithms.svelte";

export interface IMOPSOConfig extends IAlgorithmSwarmConfigBase {
    archiveSize: number,
    numberOfGrids: number,
    mutationRate: number,
    maxVelocity: number,
    c1: number,
    c2: number,
    w: number,
}

export const mopsoConfig = $state<IMOPSOConfig>({
    iterations: 300,
    population: 100,
    archiveSize: 100,
    numberOfGrids: 20,
    mutationRate: 0.5,
    maxVelocity: 5,
    c1: 2,
    c2: 2,
    w: 0.4,
    type: 'Swarm',
})
