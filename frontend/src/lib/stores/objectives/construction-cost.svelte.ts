export interface IConstructionCostConfig {
  FrequencyMatrixFilePath: string;
  DistanceMatrixFilePath: string;
  AlphaCCPenalty: number,
  GeneralQAP: boolean,
}


export const constructionCostConfig = $state<IConstructionCostConfig>({
  AlphaCCPenalty: 100,
  FrequencyMatrixFilePath: '',
  DistanceMatrixFilePath: '',
  GeneralQAP: false,
})