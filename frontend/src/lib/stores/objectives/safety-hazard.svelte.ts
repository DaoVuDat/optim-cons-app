export interface ISafetyHazardConfig {
  SEMatrixFilePath: string;
  AlphaSafetyHazardPenalty: number,
}


export const safetyHazardConfig = $state<ISafetyHazardConfig>({
  AlphaSafetyHazardPenalty: 100,
  SEMatrixFilePath: '',
})