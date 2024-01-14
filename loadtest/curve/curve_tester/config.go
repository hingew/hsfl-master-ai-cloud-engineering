package curve_tester

import (
	"encoding/json"
	"io/ioutil"
)

type NextCurvePoint struct {
	Seconds2TargetRPS int `json:"seconds2TargetRPS"`
	TargetRPS         int `json:"targetRPS"`
}

type LoadtestConfig struct {
	CurvePoints []NextCurvePoint `json:"curve"`
	Target      string           `json:"target"`
	Paths       []string         `json:"paths"`
}

func ReadConfig(filePath string) (*LoadtestConfig, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config LoadtestConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
