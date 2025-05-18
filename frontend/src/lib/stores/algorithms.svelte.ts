import {objectiveStore} from "$lib/stores/objectives.svelte";
import {
    type IGWOConfig,
    type IMOAHAConfig,
    type IAHAConfig,
    type IGAConfig,
    type IOMOAHAConfig,
    type IMOPSOConfig,
    type INSGAIIConfig,
    moahaConfig, gwoConfig, ahaConfig, gaConfig, omoahaConfig, mopsoConfig, nsgaiiConfig
} from "$lib/stores/algorithms";
import {algorithms} from "$lib/wailsjs/go/models";

export type AlgorithmConfigMap = {
    [algorithms.AlgorithmType.GeneticAlgorithm]: IGAConfig;
    [algorithms.AlgorithmType.AHA]: IAHAConfig;
    [algorithms.AlgorithmType.MOAHA]: IMOAHAConfig;
    [algorithms.AlgorithmType.GWO]: IGWOConfig;
    [algorithms.AlgorithmType.oMOAHA]: IOMOAHAConfig;
    [algorithms.AlgorithmType.MOPSO]: IMOPSOConfig;
    [algorithms.AlgorithmType.NSGAII]: INSGAIIConfig;
}

export interface AlgorithmWithLabel {
    label: string
    value: algorithms.AlgorithmType
}

const SingleList: AlgorithmWithLabel[] = [
    {
        label: 'Genetic Algorithm',
        value: algorithms.AlgorithmType.GeneticAlgorithm,
    },
    {
        label: "Artificial Hummingbird Algorithm",
        value: algorithms.AlgorithmType.AHA,
    },
    {
        label: 'Grey Wolf Algorithm',
        value: algorithms.AlgorithmType.GWO,
    },
]

const MultiList: AlgorithmWithLabel[] = [
    {
        label: 'Multi-Objective Artificial Hummingbird Algorithm',
        value: algorithms.AlgorithmType.MOAHA,
    },
    {
        label: 'OBL Multi-Objective Artificial Hummingbird Algorithm (oMOAHA)',
        value: algorithms.AlgorithmType.oMOAHA,
    },
    {
        label: 'Multi-Objective Particle Swarm Optimization (MOPSO)',
        value: algorithms.AlgorithmType.MOPSO,
    },
    {
        label: 'Non-dominated Sorting Genetic Algorithm II (NSGA-II)',
        value: algorithms.AlgorithmType.NSGAII,
    }
]

export interface IAlgorithmSwarmConfigBase {
    iterations: number
    population: number,
    type: 'Swarm'
}

export interface IAlgorithmBiologyConfigBase {
    generation: number
    chromosome: number,
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

    resetSelection = () => {
        this.selectedAlgorithm = undefined
    }

    getConfig = <T extends algorithms.AlgorithmType>(algo: T): AlgorithmConfigMap[T] => {
        switch (algo) {
            case algorithms.AlgorithmType.GeneticAlgorithm:
                return gaConfig as AlgorithmConfigMap[T]
            case algorithms.AlgorithmType.AHA:
                return ahaConfig as AlgorithmConfigMap[T]
            case algorithms.AlgorithmType.MOAHA:
                return moahaConfig as AlgorithmConfigMap[T]
            case algorithms.AlgorithmType.GWO:
                return gwoConfig as AlgorithmConfigMap[T]
            case algorithms.AlgorithmType.oMOAHA:
                return omoahaConfig as AlgorithmConfigMap[T]
            case algorithms.AlgorithmType.MOPSO:
                return mopsoConfig as AlgorithmConfigMap[T]
            case algorithms.AlgorithmType.NSGAII:
                return nsgaiiConfig as AlgorithmConfigMap[T]
        }
    }

    getValidSelection = () => {
        return this.selectedAlgorithm && this.validAlgorithmsList.find(
            a => a.value === this.selectedAlgorithm?.value
        )
    }
}

export const algorithmsStore = new AlgorithmStore()