

export enum GridFile {
  Facility,
  Phase,
}

export interface IGridConfig {
  length: number;
  width: number;
  facilitiesFilePath: {
    label: GridFile
    value: string
  };
  phasesFilePath: {
    label: GridFile,
    value: string
  };
  gridSize: number;
}


export const gridProblemConfig = $state<IGridConfig>({
  length: 120,
  width: 95,
  facilitiesFilePath: {
    label: GridFile.Facility,
    value: ''
  },
  phasesFilePath: {
    label: GridFile.Phase,
    value: ''
  },
  gridSize: 1,
})