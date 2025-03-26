package cons_lay

import "golang-moaha-construction/internal/data"

const (
	HoistingObjectiveType  = "Hoisting Objective"
	HoistingTime           = "HoistingTimeData"
	PrefabricatedLocations = "PrefabricatedLocations"
	NumberOfFloors         = "Number of Floors"
	FloorHeight            = "Floor Height"
	CraneLocations         = "Crane Locations"
	ZM                     = "ZM"
	Vuvg                   = "Vuvg"
	Vlvg                   = "Vlvg"
	Vag                    = "Vag"
	Vwg                    = "Vwg"
)

var HoistingConfigs = []*data.Config{
	{
		Name: PrefabricatedLocations,
	},
	{
		Name: NumberOfFloors,
	},
	{
		Name: HoistingTime,
	},
	{
		Name: FloorHeight,
	},
	{
		Name: CraneLocations,
	},
	{
		Name: ZM,
	},
	{
		Name: Vuvg,
	},
	{
		Name: Vlvg,
	},
	{
		Name: Vag,
	},
	{
		Name: Vwg,
	},
}

type Crane struct {
	Location
	Radius float64
}

type HoistingObjective struct {
}

func (obj *HoistingObjective) LoadHoistingData(configs []data.Config) error {

	return nil
}
