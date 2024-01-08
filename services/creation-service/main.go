package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/client"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service/pdf/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApplicationConfig struct {
	Port                 int
	TemplatingServiceURL string
}

func LoadConfigFromEnv() (*ApplicationConfig, error) {
	var config ApplicationConfig
	portStr := os.Getenv("PORT")

	if port, err := strconv.Atoi(portStr); err == nil {
		config.Port = port
	} else {

		config.Port = 3000
	}

	config.TemplatingServiceURL = os.Getenv("TEMPLATE_ENDPOINT")
	return &config, nil
}

func main() {
	config, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	grpcConn, err := grpc.Dial(config.TemplatingServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("GRPC could not connect: %v", err)
	}
	defer grpcConn.Close()

	templatingServiceClient := client.NewGrpcClient(grpcConn)

	pdf := pdf.New()
	controller := controller.NewController(pdf, templatingServiceClient)

	router := router.New()
	router.GET("/api/health/creation", health.Check)
	router.POST("/api/render/:id", controller.CreatePdf)

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
