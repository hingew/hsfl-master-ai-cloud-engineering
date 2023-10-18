package main

import (
	"fmt"
	"log"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/config"
)

type ApplicationConfig struct {
	Port int `yaml:"port"`
}

func main() {
    var appConfig *ApplicationConfig

    if err := config.Load(appConfig); err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
    }

    route := router

	addr := fmt.Sprintf("0.0.0.0:%d", *&appConfig.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}


}
