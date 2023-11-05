package auth_middleware

type Repository interface {
	VerifyToken(token string) (uint64, error)
}
