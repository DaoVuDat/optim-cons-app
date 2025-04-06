
export interface IRiskConfig {
    HazardInteractionMatrixFilePath: string;
    Delta:                   number,
    AlphaRiskPenalty:        number,
}


export const riskConfig = $state<IRiskConfig>({
    AlphaRiskPenalty: 100,
    Delta: 0.01,
    HazardInteractionMatrixFilePath: '',
})