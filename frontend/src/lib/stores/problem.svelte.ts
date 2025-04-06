import {
  continuousProblemConfig, gridProblemConfig,
  type IContinuousConfig,
  type IGridConfig,
  type IPredeterminatedConfig, predeterminedProblemConfig
} from "$lib/stores/problems";
import {objectives} from "$lib/wailsjs/go/models";

export interface ProblemWithLabel {
  label: string
  value: objectives.ProblemType
}

export const problemList: ProblemWithLabel[] = [
  {
    label: "Continuous",
    value: objectives.ProblemType.ContinuousConstructionLayout,
  },
  {
    label: "Grid",
    value: objectives.ProblemType.GridConstructionLayout,
  },
  {
    label: "Pre-determinated locations",
    value: objectives.ProblemType.PredeterminedConstructionLayout
  }
]

export type ProblemConfigMap = {
  [objectives.ProblemType.ContinuousConstructionLayout]: IContinuousConfig;
  [objectives.ProblemType.GridConstructionLayout]: IGridConfig;
  [objectives.ProblemType.PredeterminedConstructionLayout]: IPredeterminatedConfig;
}

class ProblemStore {
  selectedProblem = $state<ProblemWithLabel>()

  getConfig = <T extends objectives.ProblemType>(prob: T) : ProblemConfigMap[T]=> {
    switch (prob) {
      case objectives.ProblemType.ContinuousConstructionLayout:
        return continuousProblemConfig as ProblemConfigMap[T]
      case objectives.ProblemType.GridConstructionLayout:
        return gridProblemConfig as ProblemConfigMap[T]
      case objectives.ProblemType.PredeterminedConstructionLayout:
        return predeterminedProblemConfig as ProblemConfigMap[T]
    }
  }

  getValidSelection = () => {
    return this.selectedProblem
  }
}

export const problemStore = new ProblemStore()