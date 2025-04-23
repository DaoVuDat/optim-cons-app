export interface IPredeterminedConfig {
  numberOfLocations: number;
  numberOfFacilities: number;
}

export interface IFixedFacility {
  LocName: string
  FacilityName: string
}


class PredeterminedProblemStore {
  value = $state<IPredeterminedConfig>({
    numberOfLocations: 0,
    numberOfFacilities: 0,
  })

  locationNames = $derived.by(() =>
    Array.from({length: this.value.numberOfLocations ?? 0}, (_, i) => `L${i + 1}`)
  );

  facilityNames = $derived.by(() =>
    Array.from({length: this.value.numberOfFacilities ?? 0}, (_, i) => `TF${i + 1}`)
  );

  fixedFacilities = $state<IFixedFacility[]>([])

  setupFixedFacilities = (toBeSavedFixedFacilities: IFixedFacility[]) => {
    // this.fixedFacilities = this.fixedFacilities.splice(0, toBeSavedFixedFacilities.length, ...toBeSavedFixedFacilities);
    this.fixedFacilities = toBeSavedFixedFacilities
  }
}

export const predeterminedProblemConfig = new PredeterminedProblemStore()