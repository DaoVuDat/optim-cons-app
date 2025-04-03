
// simulate enum algorithm name
enum AlgorithmSwarmType{
  GWO,
  AHA,
  MOAHA,
}

enum AlgorithmBiologyType {
  GA
}

interface IAlgorithmSwarmConfigBase {
  name?: AlgorithmSwarmType;
  iterations: number
  population: number
}

interface IAlgorithmBiologyConfigBase {
  name?: AlgorithmBiologyType;
  generation: number
  population: number
}

interface IGWOConfig extends IAlgorithmSwarmConfigBase {}

interface IAHAConfig extends IAlgorithmSwarmConfigBase{}

interface IMOAHAConfig extends IAlgorithmSwarmConfigBase{

}

interface IGAConfig extends IAlgorithmBiologyConfigBase{
  mutationRate: number
  crossoverRate: number
}

interface IAlgorithm {
  name?: AlgorithmSwarmType | AlgorithmBiologyType
  config?: IGWOConfig | IAHAConfig | IMOAHAConfig | IGAConfig
}

class AlgorithmStore {
  algorithm = $state<IAlgorithm>({})
}

export const algorithmsStore = new AlgorithmStore()