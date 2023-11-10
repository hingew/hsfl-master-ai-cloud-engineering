package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/database"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/api/handler"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `yaml:"database"`
	Jwt      auth.JwtConfig      `yaml:"jwt"`
	Port     int                 `yaml:"port"`
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

	appConfig, err := LoadConfigFromFile(*configPath)
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
	router.POST("/auth/register", handler.Register(userRepository, hasher))
	router.POST("/auth/login", handler.Login(userRepository, hasher, tokenGenerator))

	addr := fmt.Sprintf("0.0.0.0:%d", *&appConfig.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
