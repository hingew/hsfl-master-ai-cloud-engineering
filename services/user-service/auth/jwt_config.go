package auth

import (
	"crypto/x509"
	"encoding/pem"
	"os"
)

type JwtConfig struct {
	PrivateKey string
}

func LoadConfigFromEnv() (JwtConfig) {
	return JwtConfig{PrivateKey: os.Getenv("AUTH_SIGN_KEY")}
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	block, _ := pem.Decode([]byte(config.PrivateKey))
	return x509.ParseECPrivateKey(block.Bytes)
}
