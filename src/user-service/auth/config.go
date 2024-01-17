package auth

import "time"

type Config interface {
	ReadPrivateKey() (any, error)
	ReadPublicKey() (any, error)
	ReadExpiration() (time.Duration, error)
}
