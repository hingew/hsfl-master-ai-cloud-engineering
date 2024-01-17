package ramp_tester

import (
	"encoding/json"
	"io/ioutil"
)

type RampSpecification struct {
	Duration  int `json:"duration"`
	TargetRPS int `json:"targetRPS"`
}

type TesterConfig struct {
	RampSpecifications []RampSpecification `json:"rampSpecifications"`
	Target             string              `json:"target"`
	Paths              []string            `json:"paths"`
}

func ReadConfig(filePath string) (*TesterConfig, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config TesterConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
