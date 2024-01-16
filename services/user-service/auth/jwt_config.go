package auth

import (
	"crypto/ecdsa"
	"encoding/pem"
	"os"

	"github.com/golang-jwt/jwt"
)

type JwtConfig struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
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

func readPrivateKey(stringKey string) (*ecdsa.PrivateKey, error) {
	key, _ := pem.Decode([]byte(stringKey))
	return jwt.ParseECPrivateKeyFromPEM(key.Bytes)
}

func readPublicKey(stringKey string) (*ecdsa.PublicKey, error) {
	key, _ := pem.Decode([]byte(stringKey))
	return jwt.ParseECPublicKeyFromPEM(key.Bytes)
}
