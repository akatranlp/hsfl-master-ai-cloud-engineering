package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenGenerator struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	tokenExp   time.Duration
}

func NewJwtTokenGenerator(config Config) (*JwtTokenGenerator, error) {
	uncheckedPrivateKey, err := config.ReadPrivateKey()
	if err != nil {
		return nil, err
	}
	uncheckedPublicKey, err := config.ReadPublicKey()
	if err != nil {
		return nil, err
	}
	expiration, err := config.ReadExpiration()
	if err != nil {
		return nil, err
	}

	privateKey, ok := uncheckedPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not of type *rsa.PrivateKey")
	}
	publicKey, ok := uncheckedPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not of type *rsa.PublicKey")
	}

	return &JwtTokenGenerator{privateKey, publicKey, expiration}, nil
}

func (gen *JwtTokenGenerator) CreateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}

	now := time.Now()

	jwtClaims["iat"] = now.Unix()
	jwtClaims["nbf"] = now.Unix()
	jwtClaims["exp"] = now.Add(gen.tokenExp).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)
	return token.SignedString(gen.privateKey)
}

func (gen *JwtTokenGenerator) VerifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return gen.publicKey, nil
		},
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (gen *JwtTokenGenerator) GetTokenExpiration() time.Duration {
	return gen.tokenExp
}
