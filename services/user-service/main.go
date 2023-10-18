package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/config"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/api/handler"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/api/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `yaml:"database"`
	Jwt      auth.JwtConfig      `yaml:"jwt"`
	Port     int                 `yaml:"port"`
}

func main() {
	var appConfig *ApplicationConfig

	if err := config.Load(appConfig); err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	userRepository, err := user.NewPsqlRepository(appConfig.Database)
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

	router := router.New(
		handler.NewRegisterHandler(userRepository, hasher),
		handler.NewLoginHandler(userRepository, hasher, tokenGenerator),
	)

	addr := fmt.Sprintf("0.0.0.0:%d", *&appConfig.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
