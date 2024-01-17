package main

import (
	"log"
	"net"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/proto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/api/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/controller"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/repository"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user/server"
	"google.golang.org/grpc"
)

type ApplicationConfig struct {
	Database database.PsqlConfig
	Jwt      auth.JwtConfig
}

func LoadConfigFromEnv() (*ApplicationConfig, error) {
	var err error

	databaseConfig, err := database.LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	jwtConfig, err := auth.LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	config := ApplicationConfig{
		Database: *databaseConfig,
		Jwt:      *jwtConfig,
	}

	return &config, nil
}

func main() {
	appConfig, err := LoadConfigFromEnv()
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	userRepository, err := repository.NewPsqlRepository(appConfig.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}

	if err := userRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	tokenGenerator, err := auth.NewJwtTokenGenerator(appConfig.Jwt)
	if err != nil {
		log.Fatalf("could not create JWT token generator: %s", err.Error())
	}

	hasher := crypto.NewBcryptHasher()

	ctr := controller.NewController(userRepository, hasher, tokenGenerator)
	grpcSrv := server.NewGrpcServer(userRepository, tokenGenerator)
	router := router.NewUserRouter(ctr)

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
	proto.RegisterUserServiceServer(srv, grpcSrv)

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("GRPC could not serve: %v", err)
	}

}
