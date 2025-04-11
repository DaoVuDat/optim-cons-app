
export interface IRiskConfig {
    HazardInteractionMatrixFilePath: string;
    Delta:                   number,
    AlphaRiskPenalty:        number,
}


export const riskConfig = $state<IRiskConfig>({
    AlphaRiskPenalty: 100,
    Delta: 0.01,
    HazardInteractionMatrixFilePath: '/home/daovudat/Downloads/data/conslay/f2_risk_data.xlsx',
})