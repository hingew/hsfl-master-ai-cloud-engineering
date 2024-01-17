package auth

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt"
)

type JwtConfig struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func LoadConfigFromEnv() (*JwtConfig, error) {
	privateKey, err := readPrivateKey(os.Getenv("AUTH_PRIVATE_KEY"))
	if err != nil {
		return nil, err
	}

	publicKey, err := readPublicKey(os.Getenv("AUTH_PUBLIC_KEY"))
	if err != nil {
		return nil, err
	}

	return &JwtConfig{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil

}

func readPrivateKey(stringKey string) (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(stringKey))
}

func readPublicKey(stringKey string) (*rsa.PublicKey, error) {
	return jwt.ParseRSAPublicKeyFromPEM([]byte(stringKey))
}
