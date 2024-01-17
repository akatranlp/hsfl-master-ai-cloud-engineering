package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	PrivateKeyPath         string        `env:"PRIVATE_KEY_PATH"`
	PrivateKey             string        `env:"PRIVATE_KEY"`
	PublicKeyPath          string        `env:"PUBLIC_KEY_PATH"`
	PublicKey              string        `env:"PUBLIC_KEY"`
	AccessTokenExpiration  time.Duration `env:"ACCESS_TOKEN_EXPIRATION" envDefault:"15m"`
	RefreshTokenExpiration time.Duration `env:"REFRESH_TOKEN_EXPIRATION" envDefault:"168h"`
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
