package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH"`
	PrivateKey     string `env:"PRIVATE_KEY"`
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
	PublicKey      string `env:"PUBLIC_KEY"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	var bytes []byte
	if config.PrivateKey != "" {
		bytes = []byte(config.PrivateKey)
	} else {
		bytes1, err := os.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
		bytes = bytes1
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (config JwtConfig) ReadPublicKey() (any, error) {
	var bytes []byte
	if config.PublicKey != "" {
		bytes = []byte(config.PublicKey)
	} else {
		bytes1, err := os.ReadFile(config.PublicKeyPath)
		if err != nil {
			return nil, err
		}
		bytes = bytes1
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
