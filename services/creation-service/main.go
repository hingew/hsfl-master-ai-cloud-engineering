package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/config"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
)

type ApplicationConfig struct {
	Port int `yaml:"port"`
}

func main() {
	var appConfig *ApplicationConfig

	if err := config.Load(appConfig); err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	pdf := pdf.New()
	controller := controller.NewController(pdf)

	router := router.New()
    router.POST("/render/:id", controller.CreatePdf)

	addr := fmt.Sprintf("0.0.0.0:%d", *&appConfig.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}

}
