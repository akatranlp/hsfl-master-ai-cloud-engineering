package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJwtAuthorizer(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	publicKey := &privateKey.PublicKey
	tokenGenerator := JwtTokenGenerator{privateKey, publicKey}

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
		t.Run("should error if token is not valid", func(t *testing.T) {
			// given
			token := "invalid token"

			// when
			claims, err := tokenGenerator.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Nil(t, claims)
		})

		t.Run("should error if token is not signed with RSA", func(t *testing.T) {
			// given
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"name": "Toni Tester", "exp": time.Now().Add(1 * time.Hour).Unix()})
			token, _ := jwtToken.SignedString([]byte("secret"))

			// when
			claims, err := tokenGenerator.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Nil(t, claims)
		})

		t.Run("should error if token is not signed with private key", func(t *testing.T) {
			// given
			otherPrivateKey, _ := rsa.GenerateKey(rand.Reader, 4096)

			jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"name": "Toni Tester", "exp": time.Now().Add(1 * time.Hour).Unix()})
			token, _ := jwtToken.SignedString(otherPrivateKey)

			// when
			claims, err := tokenGenerator.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Nil(t, claims)
		})

		t.Run("should error if token is expired", func(t *testing.T) {
			// given
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"name": "Toni Tester", "exp": 12345})
			token, _ := jwtToken.SignedString(privateKey)

			// when
			claims, err := tokenGenerator.VerifyToken(token)

			// then
			assert.Error(t, err)
			assert.Nil(t, claims)
		})

		t.Run("should return claims if token is valid", func(t *testing.T) {
			// given
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"name": "Toni Tester", "exp": time.Now().Add(1 * time.Hour).Unix()})
			token, _ := jwtToken.SignedString(privateKey)

			// when
			claims, err := tokenGenerator.VerifyToken(token)

			// then
			assert.NoError(t, err)
			assert.Equal(t, "Toni Tester", claims["name"])
		})
	})
}
