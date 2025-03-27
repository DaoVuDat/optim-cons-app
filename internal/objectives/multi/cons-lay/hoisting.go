package cons_lay

const (
	HoistingObjectiveType  = "Hoisting Objective"
	HoistingTimeData       = "HoistingTimeData"
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

type HoistingConfigs struct {
}

type Crane struct {
	Location
	Radius float64
}

type HoistingTime struct {
	Coordinate     Coordinate
	HoistingNumber int
	Name           string
	BuildingName   string
}

type HoistingObjective struct {
	PrefabricatedLocations []Location
	NumberOfFloors         int
	HoistingTime           []HoistingTime
	FloorHeight            float64
	CraneLocations         []Crane
	ZM                     float64
	Vuvg                   float64
	Vlvg                   float64
	Vag                    float64
	Vwg                    float64
}

func (obj *HoistingObjective) Eval() float64 {
	return 0
}

func CoverRangeOfCrane(crane Crane, buildings []Location) bool {

	return true
}
