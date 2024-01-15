package database

import (
	"fmt"
	"os"
	"strconv"
)

type PsqlConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func LoadConfigFromEnv() (*PsqlConfig, error) {
	var err error

	config := PsqlConfig{}
	config.Host = os.Getenv("POSTGRES_HOST")
	config.Port, err = strconv.Atoi(os.Getenv("POSTGRES_PORT"))

	if err != nil {
		return nil, err
	}

	config.Username = os.Getenv("POSTGRES_USERNAME")
	config.Password = os.Getenv("POSTGRES_PASSWORD")
	config.Database = os.Getenv("POSTGRES_DBNAME")

	return &config, nil
}

func (config PsqlConfig) Dsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.Database)
}
