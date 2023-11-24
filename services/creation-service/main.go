package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/client"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Port                 int    `yaml:"port"`
	TemplatingServiceURL string `yaml:"templating_service_url"`
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

	config, err := LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	templatingServiceClient := client.NewClient(config.TemplatingServiceURL)

	pdf := pdf.New()
	controller := controller.NewController(pdf, &templatingServiceClient)

	router := router.New()
    router.GET("/api/health/creation", health.Check)
	router.POST("/api/render/:id", controller.CreatePdf)

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}

}
