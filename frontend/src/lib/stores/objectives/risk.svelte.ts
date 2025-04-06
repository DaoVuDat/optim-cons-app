
export interface IGridConfig {
  length: number;
  width: number;
  facilitiesFilePath: string;
  phasesFilePath: string;
  gridSize: string;
}


export const gridProblemConfig = $state<IGridConfig>({
  length: 120,
  width: 95,
  facilitiesFilePath: '',
  phasesFilePath: '',
  gridSize: '2x2',
})