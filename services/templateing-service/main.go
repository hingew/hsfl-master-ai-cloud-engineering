package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/api/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `yaml:"database"`
}

func LoadTestData(path string) (*[]model.PdfTemplate, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var testdata []model.PdfTemplate
	if err := json.NewDecoder(f).Decode(&testdata); err != nil {
		return nil, err
	}

	return &testdata, nil
}

func LoadConfigFromFile(path string) (*ApplicationConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var config ApplicationConfig
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	use_testdata := os.Getenv("USE_TESTDATA")

	var testdata []model.PdfTemplate
	var err error

	if use_testdata == "true" {
		log.Print("Use testdata")
		p, err := LoadTestData("test_data.json")
		if err != nil {
			log.Fatalf("could not load testdata: %s", err.Error())
		} else {
			testdata = *p
		}

	}

	configPath := flag.String("config", "config.yml", "The path to the configuration file")
	flag.Parse()

	config, err := LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	repo, err := repository.NewGormPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create repository: %s", err.Error())
	}

	ctr := controller.NewController(repo)
	handler := router.NewTemplateRouter(ctr)

	if err := repo.Setup(testdata); err != nil {
		log.Fatalf("could not setup database: %s", err.Error())
	}

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
