package auth

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
)

type JwtTokenGenerator struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewJwtTokenGenerator(config JwtConfig) (*JwtTokenGenerator, error) {

	return &JwtTokenGenerator{config.PrivateKey, config.PublicKey}, nil
}

func (gen *JwtTokenGenerator) CreateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwtClaims)
	return token.SignedString(gen.privateKey)
}

func (gen *JwtTokenGenerator) VerifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodES256.Name {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}

		return gen.publicKey, nil
	})

	if err != nil {
		log.Println("Error parsing jwt token ", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
