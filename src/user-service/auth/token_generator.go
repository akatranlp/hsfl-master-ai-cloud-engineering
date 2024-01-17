package auth

import "time"

type TokenGenerator interface {
	VerifyToken(tokenString string) (map[string]interface{}, error)
	CreateToken(claims map[string]interface{}) (string, error)
	GetTokenExpiration() time.Duration
}
