package load

import (
	"encoding/json"
	"io/ioutil"
)

type TesterConfig struct {
	NumberUsers int      `json:"users"`
	Rampup      int      `json:"rampup"`
	Duration    int      `json:"duration"`
	Cooldown    int      `json:"cooldown"`
	Target      string   `json:"target"`
	Path        string   `json:"path"`
	Targets     []string `json:"targets"`
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
