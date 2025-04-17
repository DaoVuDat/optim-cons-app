export interface ITransportCostConfig {
    InteractionMatrixFilePath: string;
    AlphaTCPenalty:        number,
}


export const transportCostConfig = $state<ITransportCostConfig>({
    AlphaTCPenalty: 100,
    InteractionMatrixFilePath: '/home/daovudat/Downloads/data/conslay/transport_cost_data.xlsx',
})