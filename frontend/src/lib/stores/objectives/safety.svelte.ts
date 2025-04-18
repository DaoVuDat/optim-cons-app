export interface ISafetyConfig {
    SafetyProximityMatrixFilePath: string;
    AlphaSafetyPenalty:        number,
}


export const safetyConfig = $state<ISafetyConfig>({
    AlphaSafetyPenalty: 100,
    SafetyProximityMatrixFilePath: '',
})