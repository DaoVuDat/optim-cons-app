
export interface IPredeterminedConfig {
  numberOfLocations: number;
  numberOfFacilities: number;
}


export const predeterminedProblemConfig = $state<IPredeterminedConfig>({
  numberOfLocations: 0,
  numberOfFacilities: 0,
})