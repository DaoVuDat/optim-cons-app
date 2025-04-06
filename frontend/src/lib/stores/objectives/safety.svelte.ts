
export interface IPredeterminedConfig {
  predeterminedLocationsFilePath: string;
  facilitiesFilePath: string;
}


export const predeterminedProblemConfig = $state<IPredeterminedConfig>({
  facilitiesFilePath: '',
  predeterminedLocationsFilePath: '',
})