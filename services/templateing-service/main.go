package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/middleware"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/proto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/api/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApplicationConfig struct {
	Database                database.PsqlConfig
	UserServiceGrpcEndpoint string
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

func loadConfigFromEnv() (*ApplicationConfig, error) {
	databaseConfig, err := database.LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	userServiceGrpcEndpoint := os.Getenv("USER_GRPC_ENDPOINT")

	return &ApplicationConfig{*databaseConfig, userServiceGrpcEndpoint}, nil
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

	config, err := loadConfigFromEnv()
	if err != nil {
		log.Fatal("Could not read environment variables: ", err)
	}

	repo, err := repository.NewGormPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create repository: %s", err.Error())
	}

	grpcConn, err := grpc.Dial(config.UserServiceGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("GRPC could not connect : %v", err)
	}

	authMiddleware := middleware.NewGrpcAuthMiddleware(grpcConn)

	ctr := controller.NewController(repo)
	grpcSrv := server.NewGrpcServer(repo)
	router := router.NewTemplateRouter(ctr, authMiddleware)

	if err := repo.Setup(testdata); err != nil {
		log.Fatalf("could not setup database: %s", err.Error())
	}

	go func() {
		if err := http.ListenAndServe(":3000", router); err != nil {
			log.Fatalf("error while listen and serve: %s", err.Error())
		}
	}()

	// GRPC Server
	listener, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatalf("GRPC could not listen: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterTemplateServiceServer(srv, grpcSrv)

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("GRPC could not serve: %v", err)
	}
}
