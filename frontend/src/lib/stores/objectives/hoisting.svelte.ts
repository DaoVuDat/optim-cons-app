
export interface ISelectedCrane {
    Name: string;
    HoistingTimeFilePath: string;
    ForBuilding: string;
}

export interface ISelectedCraneWithId extends ISelectedCrane {
    Id: string
}

export interface Building {
    NumberOfFloors: number;
    FloorHeight: number;
    Name: string;
}

export interface IHoistingConfig {
    CraneLocations: ISelectedCraneWithId[];
    Buildings : Building[];
    ZM: number;
    Vuvg: number;
    Vlvg: number;
    Vag: number;
    Vwg: number;
    AlphaHoistingPenalty: number;
    AlphaHoisting: number;
    BetaHoisting: number;
}


export const hoistingConfig = $state<IHoistingConfig>({
    CraneLocations: [],
    Buildings: [],
    ZM: 2,
    Vuvg: 37.5,
    Vlvg: 37.5 / 2,
    Vag: 50,
    Vwg: 0.5,
    AlphaHoistingPenalty: 1,
    AlphaHoisting: 0.25,
    BetaHoisting: 1,
})