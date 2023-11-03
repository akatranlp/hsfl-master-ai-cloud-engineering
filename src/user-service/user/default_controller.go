package user

import (
	"encoding/json"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"
)

type contextKey int

const authenticatedUserKey contextKey = 0

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (r *loginRequest) isValid() bool {
	return r.Email != "" && r.Password != ""
}

type DefaultController struct {
	userRepository Repository
	hasher         crypto.Hasher
	tokenGenerator auth.TokenGenerator
}

func NewDefaultController(
	userRepository Repository,
	hasher crypto.Hasher,
	tokenGenerator auth.TokenGenerator,
) *DefaultController {
	return &DefaultController{userRepository, hasher, tokenGenerator}
}

func (ctrl *DefaultController) Login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := ctrl.userRepository.FindByEmail(request.Email)
	if err != nil {
		log.Printf("could not find user by email: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) < 1 {
		w.Header().Add("WWW-Authenticate", "Basic realm=Restricted")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if ok := ctrl.hasher.Validate([]byte(request.Password), users[0].Password); !ok {
		w.Header().Add("WWW-Authenticate", "Basic realm=Restricted")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiration := 1 * time.Hour
	accessToken, err := ctrl.tokenGenerator.CreateToken(map[string]interface{}{
		"email": request.Email,
		"exp":   time.Now().Add(expiration).Unix(),
	})

	json.NewEncoder(w).Encode(loginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(expiration.Seconds()),
	})
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Username    string `json:"username"`
	ProfileName string `json:"profileName"`
}

func (r *registerRequest) isValid() bool {
	return r.Email != "" && r.Password != "" && r.Username != "" && r.ProfileName != ""
}

func (ctrl *DefaultController) Register(w http.ResponseWriter, r *http.Request) {
	var request registerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := ctrl.userRepository.FindByEmail(request.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(user) > 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}

	hashedPassword, err := ctrl.hasher.Hash([]byte(request.Password))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := ctrl.userRepository.Create([]*model.DbUser{{
		Email:       request.Email,
		Password:    hashedPassword,
		Username:    request.Username,
		ProfileName: request.ProfileName,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ctrl.userRepository.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (ctrl *DefaultController) GetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authenticatedUserKey).(model.DbUser)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (ctrl *DefaultController) GetUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	user, err := ctrl.userRepository.FindByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (ctrl *DefaultController) PutUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	var request model.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ctrl.userRepository.Update(username, &request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	if err := ctrl.userRepository.Delete([]*model.DbUser{{Username: username}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) AuthenticationMiddleWare(w http.ResponseWriter, r *http.Request, next router.Next) {
	token := r.Header.Get("Authorization")

	after, found := strings.CutPrefix(token, "Bearer ")
	if !found {
		http.Error(w, "There was no Token provided", http.StatusUnauthorized)
		return
	}

	claims, err := ctrl.tokenGenerator.VerifyToken(after)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)

	if !ok {
		http.Error(w, "There is no email claim in your token", http.StatusUnauthorized)
		return
	}

	users, err := ctrl.userRepository.FindByEmail(email)
	if err != nil {
		http.Error(w, "The user doesn't exist anymore", http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), authenticatedUserKey, *users[0])
	next(r.WithContext(ctx))
}
