package auth_middleware

import (
	"context"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"net/http"
)

type contextKey string

const (
	AuthenticatedUserId contextKey = "user"
)

type DefaultController struct {
	authRepository Repository
}

func NewDefaultController(
	authRepository Repository,
) *DefaultController {
	return &DefaultController{authRepository}
}

func (ctrl *DefaultController) AuthenticationMiddleware(w http.ResponseWriter, r *http.Request, next router.Next) {
	ctx := context.WithValue(r.Context(), AuthenticatedUserId, uint64(1))
	next(r.WithContext(ctx))
	// TODO: Reactivate if we shall use Authentication
	/*
		bearerToken := r.Header.Get("Authorization")
		token, found := strings.CutPrefix(bearerToken, "Bearer ")
		if !found {
			http.Error(w, "There was no Token provided", http.StatusUnauthorized)
			return
		}

		userId, err := ctrl.authRepository.VerifyToken(token)
		if err != nil {
			http.Error(w, "There was an Error while verifying you token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authenticatedUserId, userId)
		next(r.WithContext(ctx))
	*/
}
