package objectives

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"log"
	"testing"
)

func TestTransportCostObjectiveMini_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]data.Location
		expected  float64
		name      string
	}{
		{
			locations: CreateInputMini(),
			expected:  5397.28,
			name:      "mini locations",
		},
	}

	interactionMatrix, err := ReadInteractionTransportCostDataFromFile("../../../data/conslay/mini/transport_cost_data.xlsx")

	transportCostCfg := TransportCostConfigs{
		InteractionMatrix: interactionMatrix,
		AlphaTCPenalty:    100,
		Phases:            CreateInputPhasesMini(),
	}

	tcObj, err := CreateTransportCostObjectiveFromConfig(transportCostCfg)
	if err != nil {
		log.Fatal(err)
	}

	// calculate result
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			result := tcObj.Eval(test.locations)
			if util.RoundTo(result, 2) != test.expected {
				t.Errorf("expected result to be %f, got %f", test.expected, result)
			}

		})
	}
}
