package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenGenerator struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
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

	privateKey, ok := uncheckedPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, err
	}
	publicKey, ok := uncheckedPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	return &JwtTokenGenerator{privateKey, publicKey}, nil
}

func (gen *JwtTokenGenerator) CreateToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}

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
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
