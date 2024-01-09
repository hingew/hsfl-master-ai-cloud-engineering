package main

import (
	"log"
	"net/http"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/api/handler"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user"
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

	config := ApplicationConfig{
		Database: *databaseConfig,
		Jwt:      auth.LoadConfigFromEnv(),
	}

	return &config, nil
}

func main() {
	appConfig, err := LoadConfigFromEnv()
	if err != nil {
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

	router := router.New()
	router.GET("/api/health/user", health.Check)

	router.POST("/auth/register", handler.Register(userRepository, hasher))
	router.POST("/auth/login", handler.Login(userRepository, hasher, tokenGenerator))

	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
