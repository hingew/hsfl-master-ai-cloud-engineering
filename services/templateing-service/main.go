package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/api/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `yaml:"database"`
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
	configPath := flag.String("config", "config.yml", "The path to the configuration file")
	flag.Parse()

	// TODO: Haukes lib Funktion hier nutzen
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

	if err := repo.Setup(); err != nil {
		log.Fatalf("could not setup database: %s", err.Error())
	}

	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
