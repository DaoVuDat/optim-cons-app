enum ProblemType {
  Grid,
  Continuous,
  PreLocated
}

interface IGridConfig {}

interface IContinuousConfig {}

interface IPreLocatedConfig {}

interface IProblem {
  name?: ProblemType
  config? :IContinuousConfig | IGridConfig | IPreLocatedConfig
}

class ProblemStore {
  problem = $state<IProblem>({})
}

export const problemStore = new ProblemStore()