package auth

import (
	"crypto/ecdsa"

	"github.com/golang-jwt/jwt"
)

type JwtTokenGenerator struct {
	privateKey *ecdsa.PrivateKey
}

func NewJwtTokenGenerator(config Config) (*JwtTokenGenerator, error) {
	key, err := config.ReadPrivateKey()
	if err != nil {
		return nil, err
	}

	privateKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, err
	}

	return &JwtTokenGenerator{privateKey}, nil
}

func (gen *JwtTokenGenerator) CreateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims)
	return token.SignedString(gen.privateKey)
}
