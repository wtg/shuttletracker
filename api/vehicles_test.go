package api

import "testing"

func TestCardinalDirection(t *testing.T) {
	table := [][]string{
		{"0", "North"},
		{"45", "North-East"},
		{"90", "East"},
		{"135", "South-East"},
		{"180", "South"},
		{"225", "South-West"},
		{"270", "West"},
		{"315", "North-West"},
		{"this isn't a direction lol", "North"},
	}

	for _, testCase := range table {
		direction := CardinalDirection(&testCase[0])
		expected := testCase[1]
		if direction != expected {
			t.Errorf("Got %v, expected %v.", direction, expected)
		}
	}
}
