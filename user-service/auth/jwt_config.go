package auth

import (
	"crypto/x509"
	"encoding/pem"
	"os"
)

type JwtConfig struct {
	SignKey string `yaml:"signKey"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	bytes, err := os.ReadFile(config.SignKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	return x509.ParseECPrivateKey(block.Bytes)
}
