export enum ContinuousFile {
  Facility,
  Phase,
}

export interface IContinuousConfig {
  length: number;
  width: number;
  facilitiesFilePath: {
    label: ContinuousFile
    value: string
  };
  phasesFilePath: {
    label: ContinuousFile,
    value: string
  }
}


export const continuousProblemConfig = $state<IContinuousConfig>({
  length: 120,
  width: 95,
  facilitiesFilePath: {
    label: ContinuousFile.Facility,
    value: '/home/daovudat/Downloads/data/conslay/locations.xlsx'
  },
  phasesFilePath: {
    label: ContinuousFile.Phase,
    value: '/home/daovudat/Downloads/data/conslay/dynamicBuilding.xlsx'
  },
})