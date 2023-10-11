package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/api/handler"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/api/router"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/user-service/user"
	"gopkg.in/yaml.v3"
)

type ApplicationConfig struct {
	Database database.PsqlConfig `yaml:"database"`
	Jwt      auth.JwtConfig      `yaml:"jwt"`
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
	port := flag.String("port", "8080", "The listening port")
	configPath := flag.String("config", "config.yml", "The path to the configuration file")
	flag.Parse()

	config, err := LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("could not load application configuration: %s", err.Error())
	}

	userRepository, err := user.NewPsqlRepository(config.Database)
	if err != nil {
		log.Fatalf("could not create user repository: %s", err.Error())
	}

	if err := userRepository.Migrate(); err != nil {
		log.Fatalf("could not migrate: %s", err.Error())
	}

	tokenGenerator, err := auth.NewJwtTokenGenerator(config.Jwt)
	if err != nil {
		log.Fatalf("could not create JWT token generator: %s", err.Error())
	}

	hasher := crypto.NewBcryptHasher()

	handler := router.New(
		handler.NewRegisterHandler(userRepository, hasher),
		handler.NewLoginHandler(userRepository, hasher, tokenGenerator),
	)

	addr := fmt.Sprintf("0.0.0.0:%s", *port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
