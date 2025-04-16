export interface ISafetyConfig {
    SafetyProximityMatrixFilePath: string;
    AlphaSafetyPenalty:        number,
}


export const safetyConfig = $state<ISafetyConfig>({
    AlphaSafetyPenalty: 100,
    SafetyProximityMatrixFilePath: '/home/daovudat/Downloads/data/conslay/safety_data.xlsx',
})