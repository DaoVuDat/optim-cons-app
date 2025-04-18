export interface ITransportCostConfig {
  InteractionMatrixFilePath: string;
  AlphaTCPenalty: number,
}


export const transportCostConfig = $state<ITransportCostConfig>({
  AlphaTCPenalty: 100,
  InteractionMatrixFilePath: '',
})