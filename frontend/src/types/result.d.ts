import type {Facility} from "$lib/stores/problems/problem";
import type {ISelectedCrane} from "$lib/stores/objectives";

export interface Penalty {
  [k: string]: number
}

export type ValuesWithKey = { [key: string]: number }

export interface MapLocation {
  [k: string]: Facility
}

interface Crane extends Facility {
  BuildingName: string[]
  Radius: number
  CraneSymbol: string
}

export interface ResultLocation {
  MapLocations: MapLocation
  Value: number[]
  ValuesWithKey: ValuesWithKey
  Convergence: number[]
  Penalty: Penalty
  Cranes: Crane[]
  Phases: string[][]
}

export interface ResultLocationWithId extends ResultLocation {
  Id: string
}