package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	PrivateKeyPath  string        `env:"PRIVATE_KEY_PATH"`
	PrivateKey      string        `env:"PRIVATE_KEY"`
	PublicKeyPath   string        `env:"PUBLIC_KEY_PATH"`
	PublicKey       string        `env:"PUBLIC_KEY"`
	TokenExpiration time.Duration `env:"TOKEN_EXPIRATION,notEmpty"`
}

func (config JwtConfig) ReadPrivateKey() (any, error) {
	var bytes []byte
	if config.PrivateKey != "" {
		bytes = []byte(config.PrivateKey)
	} else {
		var err error
		bytes, err = os.ReadFile(config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
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
		var err error
		bytes, err = os.ReadFile(config.PublicKeyPath)
		if err != nil {
			return nil, err
		}
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func (config JwtConfig) ReadExpiration() (time.Duration, error) {
	return config.TokenExpiration, nil
}
