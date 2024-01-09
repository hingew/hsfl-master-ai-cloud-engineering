package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/model"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/proto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/api/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/templateing-service/templates/repository"
	"google.golang.org/grpc"
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

	config := ApplicationConfig{}
	config.Database.Host = os.Getenv("POSTGRES_HOST")
	config.Database.Port, err = strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatalf("could parse postgres port")
	}
	config.Database.Username = os.Getenv("POSTGRES_USERNAME")
	config.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	config.Database.Database = os.Getenv("POSTGRES_DBNAME")

	repo, err := repository.NewGormPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create repository: %s", err.Error())
	}

	ctr := controller.NewController(repo)
	grpcSrv := controller.NewGrpcServer(repo)
	handler := router.NewTemplateRouter(ctr)

	if err := repo.Setup(testdata); err != nil {
		log.Fatalf("could not setup database: %s", err.Error())
	}

	go func() {
		if err := http.ListenAndServe(":3000", handler); err != nil {
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
