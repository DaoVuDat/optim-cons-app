import {objectiveStore} from "$lib/stores/objectives.svelte";
import type {IGWOConfig, IMOAHAConfig, IAHAConfig, IGAConfig} from "$lib/stores/algorithms";

export enum Algorithms {
  GWO = "GWO",
  AHA = "AHA",
  MOAHA = "MOAHA",
  GA = "GA",
}

export interface AlgorithmWithLabel {
  label: string
  value: Algorithms
}

const SingleList: AlgorithmWithLabel[] = [
  {
    label: 'Genetic Algorithm',
    value: Algorithms.GA,
  },
  {
    label: "Artificial Hummingbird Algorithm",
    value: Algorithms.AHA,
  },
  {
    label: 'Grey Wolf Algorithm',
    value: Algorithms.GWO,
  },
]

const MultiList: AlgorithmWithLabel[] = [
  {
    label: 'Multi-Objective Artificial Hummingbird Algorithm',
    value: Algorithms.MOAHA,
  }
]

export interface IAlgorithmSwarmConfigBase {
  iterations: {
    label: "Iterations",
    value: number
  }
  population: {
    label: "Population",
    value: number
  },
  type: 'Swarm'
}

export interface IAlgorithmBiologyConfigBase {
  generation: {
    label: "Generation",
    value: number
  }
  chromosome: {
    label: "Chromosome",
    value: number
  },
  type: 'Biology'
}

class AlgorithmStore {

  validAlgorithmsList = $derived.by<AlgorithmWithLabel[]>(() => {
    if (objectiveStore.objectives.selectedObjectives.length == 1) {
      return SingleList
    } else if (objectiveStore.objectives.selectedObjectives.length >= 1) {
      return MultiList
    } else {
      return []
    }
  })

  selectedAlgorithm = $state<AlgorithmWithLabel>()

  // configs
  gwoConfig = $state<IGWOConfig>()
  ahaConfig = $state<IAHAConfig>()
  gaConfig = $state<IGAConfig>()
  moahaConfig = $state<IMOAHAConfig>()
}

export const algorithmsStore = new AlgorithmStore()