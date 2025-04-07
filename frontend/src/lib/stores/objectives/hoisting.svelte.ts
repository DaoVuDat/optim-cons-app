
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
    CraneLocations: ISelectedCrane[];
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
    CraneLocations: [],
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