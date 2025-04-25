import {
  continuousProblemConfig, gridProblemConfig,
  type IContinuousConfig,
  type IGridConfig,
  type IPredeterminedConfig, predeterminedProblemConfig
} from "$lib/stores/problems";
import {data} from "$lib/wailsjs/go/models";
import {objectiveStore} from "$lib/stores/objectives.svelte";

export interface ProblemWithLabel {
  label: string
  value: data.ProblemName
}

export const problemList: ProblemWithLabel[] = [
  {
    label: "Continuous",
    value: data.ProblemName.ContinuousConstructionLayout,
  },
  {
    label: "Grid",
    value: data.ProblemName.GridConstructionLayout,
  },
  {
    label: "Pre-determined locations",
    value: data.ProblemName.PredeterminedConstructionLayout
  }
]

export type ProblemConfigMap = {
  [data.ProblemName.ContinuousConstructionLayout]: IContinuousConfig;
  [data.ProblemName.GridConstructionLayout]: IGridConfig;
  [data.ProblemName.PredeterminedConstructionLayout]: IPredeterminedConfig;
}

class ProblemStore {
  selectedProblem = $state<ProblemWithLabel>()

  validProblemList = $derived.by<ProblemWithLabel[]>(() => {
    if (objectiveStore.objectives.selectedObjectives.find(
        o => o.objectiveType === data.ObjectiveType.ConstructionCostObjective)) {
      return problemList.filter(prob => prob.value === data.ProblemName.PredeterminedConstructionLayout)
    } else {
      return problemList.filter(prob => prob.value !== data.ProblemName.PredeterminedConstructionLayout)
    }
  })

  getConfig = <T extends data.ProblemName>(prob: T) : ProblemConfigMap[T]=> {
    switch (prob) {
      case data.ProblemName.ContinuousConstructionLayout:
        return continuousProblemConfig as ProblemConfigMap[T]
      case data.ProblemName.GridConstructionLayout:
        return gridProblemConfig as ProblemConfigMap[T]
      case data.ProblemName.PredeterminedConstructionLayout:
        return predeterminedProblemConfig.value as ProblemConfigMap[T]
    }
  }

  getValidSelection = () => {
    return this.selectedProblem
  }
}

export const problemStore = new ProblemStore()