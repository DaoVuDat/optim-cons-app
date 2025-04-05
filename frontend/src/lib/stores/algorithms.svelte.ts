import {objectiveStore} from "$lib/stores/objectives.svelte";
import {
  type IGWOConfig,
  type IMOAHAConfig,
  type IAHAConfig,
  type IGAConfig,
  moahaConfig, gwoConfig, ahaConfig, gaConfig
} from "$lib/stores/algorithms";

export enum Algorithms {
  GWO = "GWO",
  AHA = "AHA",
  MOAHA = "MOAHA",
  GA = "GA",
}

export type AlgorithmConfigMap = {
  [Algorithms.GA]: IGAConfig;
  [Algorithms.AHA]: IAHAConfig;
  [Algorithms.MOAHA]: IMOAHAConfig;
  [Algorithms.GWO]: IGWOConfig;
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

  getConfig = <T extends Algorithms>(algo: T) : AlgorithmConfigMap[T]=> {
    switch (algo) {
      case Algorithms.GA:
        return gaConfig as AlgorithmConfigMap[T]
      case Algorithms.AHA:
        return ahaConfig as AlgorithmConfigMap[T]
      case Algorithms.MOAHA:
        return moahaConfig as AlgorithmConfigMap[T]
      case Algorithms.GWO:
        return gwoConfig as AlgorithmConfigMap[T]
    }
  }

  getValidSelection = () => {
    return this.selectedAlgorithm && this.validAlgorithmsList.find(
      a => a.value === this.selectedAlgorithm?.value
    )
  }
}

export const algorithmsStore = new AlgorithmStore()