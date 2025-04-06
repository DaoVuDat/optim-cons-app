import {
  continuousProblemConfig, gridProblemConfig,
  type IContinuousConfig,
  type IGridConfig,
  type IPredeterminedConfig, predeterminedProblemConfig
} from "$lib/stores/problems";

export enum ProblemType {
  Grid,
  Continuous,
  PreLocated
}

export interface ProblemWithLabel {
  label: string
  value: ProblemType
}

export const problemList: ProblemWithLabel[] = [
  {
    label: "Continuous",
    value: ProblemType.Continuous,
  },
  {
    label: "Grid",
    value: ProblemType.Grid,
  },
  {
    label: "Pre-determined locations",
    value: ProblemType.PreLocated
  }
]

export type ProblemConfigMap = {
  [ProblemType.Continuous]: IContinuousConfig;
  [ProblemType.Grid]: IGridConfig;
  [ProblemType.PreLocated]: IPredeterminedConfig;
}

class ProblemStore {
  selectedProblem = $state<ProblemWithLabel>()

  getConfig = <T extends ProblemType>(prob: T) : ProblemConfigMap[T]=> {
    switch (prob) {
      case ProblemType.Continuous:
        return continuousProblemConfig as ProblemConfigMap[T]
      case ProblemType.Grid:
        return gridProblemConfig as ProblemConfigMap[T]
      case ProblemType.PreLocated:
        return predeterminedProblemConfig as ProblemConfigMap[T]
    }
  }

  getValidSelection = () => {
    return this.selectedProblem
  }
}

export const problemStore = new ProblemStore()