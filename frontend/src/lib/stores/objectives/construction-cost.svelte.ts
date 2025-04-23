export interface IConstructionCostConfig {
  FrequencyMatrixFilePath: string;
  DistanceMatrixFilePath: string;
  AlphaCCPenalty: number,
}


export const constructionCostConfig = $state<IConstructionCostConfig>({
  AlphaCCPenalty: 100,
  FrequencyMatrixFilePath: '',
  DistanceMatrixFilePath: ''
})