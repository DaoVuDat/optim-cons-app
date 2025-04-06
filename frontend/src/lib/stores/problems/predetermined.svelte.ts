
export enum PredeteriminatedFile {
  Predeterminated
}

export interface IPredeterminatedConfig {
  predeterminatedLocationsFilePath: {
    label: PredeteriminatedFile
    value: string
  };
}


export const predeterminedProblemConfig = $state<IPredeterminatedConfig>({
  predeterminatedLocationsFilePath: {
    label: PredeteriminatedFile.Predeterminated,
    value: ''
  },
})