
export interface ISelectedCrane {
    Name: string;
    BuildingNames: string;
    Radius: number;
    HoistingTimeFilePath: string;
}

export interface ISelectedCraneWithId extends ISelectedCrane {
    Id: string
}

export interface IHoistingConfig {
    CraneLocations: ISelectedCraneWithId[];
    NumberOfFloors: number;
    FloorHeight: number;
    ZM: number;
    Vuvg: number;
    Vlvg: number;
    Vag: number;
    Vwg: number;
    AlphaHoistingPenalty: number;
    AlphaHoisting: number;
    BetaHoisting: number;
    NHoisting: number;
}


export const hoistingConfig = $state<IHoistingConfig>({
    CraneLocations: [
        {
            BuildingNames: 'TF4 TF5 TF8 TF9 TF10',
            HoistingTimeFilePath: '/home/daovudat/Downloads/data/conslay/f1_hoisting_time_data.xlsx',
            Name: 'TF14',
            Radius: 40,
            Id: Math.random().toString()
        }
    ],
    NumberOfFloors: 10,
    FloorHeight: 3.2,
    ZM: 2,
    Vuvg: 37.5,
    Vlvg: 37.5 / 2,
    Vag: 50,
    Vwg: 0.5,
    AlphaHoistingPenalty: 1,
    AlphaHoisting: 0.25,
    BetaHoisting: 1,
    NHoisting: 1,
})