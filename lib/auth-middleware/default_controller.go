package auth_middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type contextKey string

const (
	AuthenticatedUserId contextKey = "user"
)

type DefaultController struct {
	authRepository Repository
	authIsActive   bool
}

func NewDefaultController(
	authRepository Repository,
	authIsActive bool,
) *DefaultController {
	return &DefaultController{authRepository, authIsActive}
}

func (ctrl *DefaultController) AuthenticationMiddleware(w http.ResponseWriter, r *http.Request, next router.Next) {
	if !ctrl.authIsActive {
		ctx := context.WithValue(r.Context(), AuthenticatedUserId, uint64(1))
		next(r.WithContext(ctx))
		return
	}

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

	ctx := context.WithValue(r.Context(), AuthenticatedUserId, userId)
	next(r.WithContext(ctx))
}
