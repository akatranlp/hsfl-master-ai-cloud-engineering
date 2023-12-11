package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/utils"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/crypto"
	"golang.org/x/sync/singleflight"

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
	authIsActive   bool
	g              *singleflight.Group
}

func NewDefaultController(
	userRepository Repository,
	hasher crypto.Hasher,
	tokenGenerator auth.TokenGenerator,
	authIsActive bool,
) *DefaultController {
	g := &singleflight.Group{}
	return &DefaultController{userRepository, hasher, tokenGenerator, authIsActive, g}
}

func (ctrl *DefaultController) createToken(userID uint64, email string, tokenVersion uint64, expiration time.Duration) (string, error) {
	return ctrl.tokenGenerator.CreateToken(map[string]interface{}{
		"id":            userID,
		"email":         email,
		"token_version": tokenVersion,
		"exp":           time.Now().Add(expiration).Unix(),
	})
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

	accessTokenExpiration := 1 * time.Hour
	accessToken, err := ctrl.createToken(users[0].ID, users[0].Email, users[0].TokenVersion, accessTokenExpiration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshTokenExpiration := 7 * 24 * time.Hour
	refreshToken, err := ctrl.createToken(users[0].ID, users[0].Email, users[0].TokenVersion, refreshTokenExpiration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(refreshTokenExpiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/v1/refresh-token",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &newCookie)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessTokenExpiration.Seconds()),
	})
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	ProfileName string `json:"profileName"`
}

func (r *registerRequest) isValid() bool {
	return r.Email != "" && r.Password != "" && r.ProfileName != ""
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
		ProfileName: request.ProfileName,
	}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *DefaultController) tokenVerification(token string) (*model.DbUser, int, error) {
	claims, err := ctrl.tokenGenerator.VerifyToken(token)
	if err != nil {
		log.Println("ERROR [tokenVerification - VerifyToken]: ", err.Error())
		return nil, http.StatusUnauthorized, errors.New("token couldn't be verified")
	}

	email, ok := claims["email"].(string)
	if !ok {
		log.Println("ERROR [tokenVerification - get email claim]: ", "There is no email claim in your token")
		return nil, http.StatusUnauthorized, errors.New("there is no email claim in your token")
	}

	tokenV, ok := claims["token_version"].(float64)
	if !ok {
		log.Println("ERROR [tokenVerification - get token_version claim]: ", "There is no token_version claim in your token")
		return nil, http.StatusUnauthorized, errors.New("there is no token_version claim in your token")
	}
	tokenVersion := uint64(tokenV)

	users, err := ctrl.userRepository.FindByEmail(email)
	if err != nil {
		log.Println("ERROR [tokenVerification - FindByEmail]: ", err.Error())
		return nil, http.StatusInternalServerError, errors.New("internal server error")
	}

	if len(users) < 1 {
		log.Println("ERROR [tokenVerification - len(users) < 1]: ", "Couldn't find user by email")
		return nil, http.StatusUnauthorized, errors.New("couldn't find user by email")
	}

	if users[0].TokenVersion != tokenVersion {
		log.Println("ERROR [tokenVerification - token version]: ", "The token version is not valid")
		return nil, http.StatusUnauthorized, errors.New("the token version is not valid")
	}

	return users[0], 200, nil
}

func (ctrl *DefaultController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Println("ERROR [REFRESH_TOKEN - get cookie]: ", err.Error())
		http.Error(w, "There was no cookie in the request!", http.StatusUnauthorized)
		return
	}

	user, statusCode, err := ctrl.tokenVerification(cookie.Value)
	if user == nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	accessTokenExpiration := 1 * time.Hour
	accessToken, err := ctrl.createToken(user.ID, user.Email, user.TokenVersion, accessTokenExpiration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	refreshTokenExpiration := 7 * 24 * time.Hour
	refreshToken, err := ctrl.createToken(user.ID, user.Email, user.TokenVersion, refreshTokenExpiration)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(refreshTokenExpiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/v1/refresh-token",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &newCookie)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessTokenExpiration.Seconds()),
	})
}

func (ctrl *DefaultController) Logout(w http.ResponseWriter, r *http.Request) {
	all := r.URL.Query().Get("all")

	if all != "" {
		user := r.Context().Value(authenticatedUserKey).(*model.DbUser)
		newTokenVersion := user.TokenVersion + 1
		if err := ctrl.userRepository.Update(user.ID, &model.DbUserPatch{TokenVersion: &newTokenVersion}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	newCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/v1/refresh-token",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &newCookie)
	w.WriteHeader(http.StatusOK)
}

func (ctrl *DefaultController) GetUsers(w http.ResponseWriter, _ *http.Request) {
	newUsers, err, _ := ctrl.g.Do("get-users", func() (interface{}, error) {
		return ctrl.userRepository.FindAll()
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	users := newUsers.([]*model.DbUser)

	userDto := utils.Map(users, func(user *model.DbUser) model.UserDTO {
		return user.ToDto()
	})

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDto)
}

func (ctrl *DefaultController) GetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authenticatedUserKey).(*model.DbUser)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToDto())
}

func (ctrl *DefaultController) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userid").(string)

	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser, err, _ := ctrl.g.Do(fmt.Sprintf("user-%d", id), func() (interface{}, error) {
		return ctrl.userRepository.FindById(id)
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := newUser.(*model.DbUser)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToDto())
}

type putMeRequest struct {
	Password    string `json:"password"`
	ProfileName string `json:"profileName"`
	Balance     *int64 `json:"balance"`
}

func (ctrl *DefaultController) PatchMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authenticatedUserKey).(*model.DbUser)

	var request putMeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var patchUser model.DbUserPatch

	if request.Password != "" {
		hashedPassword, err := ctrl.hasher.Hash([]byte(request.Password))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		patchUser.Password = &hashedPassword
	}
	if request.ProfileName != "" {
		patchUser.ProfileName = &request.ProfileName
	}
	if request.Balance != nil {
		patchUser.Balance = request.Balance
	}

	if err := ctrl.userRepository.Update(user.ID, &patchUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (ctrl *DefaultController) DeleteMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(authenticatedUserKey).(*model.DbUser)
	if err := ctrl.userRepository.Delete([]*model.DbUser{user}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type validateTokenRequest struct {
	Token string `json:"token"`
}

func (r *validateTokenRequest) isValid() bool {
	return r.Token != ""
}

func (ctrl *DefaultController) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var request validateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("ERROR [VALIDATE_TOKEN]: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		log.Println("ERROR [VALIDATE_TOKEN]: ", "Qis not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ctrl.authIsActive {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 1, "email": "test@test.com"})
		return
	}

	user, statusCode, err := ctrl.tokenVerification(request.Token)
	if user == nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToDto())
}

func (ctrl *DefaultController) MoveUserAmount(w http.ResponseWriter, r *http.Request) {
	// Fully implement this if we need Authentication ????
	var request shared_types.MoveBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payingUser, err := ctrl.userRepository.FindById(request.UserId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	receivingUser, err := ctrl.userRepository.FindById(request.ReceivingUserId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	payingUserBalance := payingUser.Balance - request.Amount
	receivingUserBalance := receivingUser.Balance + request.Amount

	userPatch := &model.DbUserPatch{Balance: &payingUserBalance}
	err = ctrl.userRepository.Update(payingUser.ID, userPatch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userPatch = &model.DbUserPatch{Balance: &receivingUserBalance}
	err = ctrl.userRepository.Update(receivingUser.ID, userPatch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shared_types.MoveBalanceResponse{Success: true})
}

func (ctrl *DefaultController) AuthenticationMiddleWare(w http.ResponseWriter, r *http.Request, next router.Next) {
	if !ctrl.authIsActive {
		user, err := ctrl.userRepository.FindById(1)
		if err != nil {
			http.Error(w, "The user doesn't exist anymore", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), authenticatedUserKey, user)
		next(r.WithContext(ctx))
		return
	}

	token := r.Header.Get("Authorization")

	after, found := strings.CutPrefix(token, "Bearer ")
	if !found {
		http.Error(w, "There was no Token provided", http.StatusUnauthorized)
		return
	}
	user, statusCode, err := ctrl.tokenVerification(after)
	if user == nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	ctx := context.WithValue(r.Context(), authenticatedUserKey, user)
	next(r.WithContext(ctx))
}
