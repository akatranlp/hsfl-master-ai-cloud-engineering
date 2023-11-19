package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH,notEmpty"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH,notEmpty"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	bytes, err := os.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (config JwtConfig) ReadPublicKey() (any, error) {
	bytes, err := os.ReadFile(config.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
