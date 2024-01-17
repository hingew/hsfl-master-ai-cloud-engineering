package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJwtAuthorizer(t *testing.T) {
    privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	publicKey := privateKey.PublicKey
	tokenGenerator := JwtTokenGenerator{privateKey, &publicKey}

	t.Run("CreateToken", func(t *testing.T) {
		t.Run("should generate valid JWT token", func(t *testing.T) {
			// given
			// when
			token, err := tokenGenerator.CreateToken(map[string]interface{}{
				"exp":  12345,
				"user": "test",
			})

			// then
			assert.NoError(t, err)
			tokenParts := strings.Split(token, ".")
			assert.Len(t, tokenParts, 3)

			b, _ := base64.
				StdEncoding.
				WithPadding(base64.NoPadding).
				DecodeString(tokenParts[1])

			var claims map[string]interface{}
			json.Unmarshal(b, &claims)

			assert.Equal(t, float64(12345), claims["exp"])
			assert.Equal(t, "test", claims["user"])
		})
	})

	t.Run("VerifyToken", func(t *testing.T) {
		t.Run("should fail the token is not valid", func(t *testing.T) {
			//given
			token := "invalid"

			//when
			claims, err := tokenGenerator.VerifyToken(token)

			//then
			assert.Error(t, err)
			assert.Nil(t, claims)

		})

		t.Run("should succeed when the token is valid", func(t *testing.T) {

			claims := map[string]interface{}{
				"exp":  time.Now().Add(1 * time.Hour).Unix(),
				"user": "test",
			}

			//given
			token, err := tokenGenerator.CreateToken(claims)
			assert.Nil(t, err)

			//when
			newClaims, err := tokenGenerator.VerifyToken(token)

			//then
			assert.Nil(t, err)
			assert.Equal(t, newClaims["user"], "test")

		})

	})
}
