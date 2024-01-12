package curve_tester

import (
	"encoding/json"
	"io/ioutil"
)

type NextGraphPoint struct {
	Seconds2TargetRPS int `json:"duration"`
	TargetRPS         int `json:"targetRPS"`
}

type LoadtestConfig struct {
	Graph  []NextGraphPoint `json:"graph"`
	Target string           `json:"target"`
	Paths  []string         `json:"paths"`
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
