package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/utils"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/auth"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/crypto"

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
		"id":    users[0].ID,
		"email": request.Email,
		"exp":   time.Now().Add(expiration).Unix(),
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(loginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(expiration.Seconds()),
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

func (ctrl *DefaultController) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := ctrl.userRepository.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	user, err := ctrl.userRepository.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToDto())
}

type putMeRequest struct {
	Password    string `json:"password"`
	ProfileName string `json:"profileName"`
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
	// Fully implement this if we need Authentication
	var request validateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("ERROR [VALIDATE_TOKEN]: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !request.isValid() {
		log.Println("ERROR [VALIDATE_TOKEN]: ", "is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": 1, "email": "test@test.com"})
}

func (ctrl *DefaultController) MoveUserAmount(w http.ResponseWriter, r *http.Request) {
	// Fully implement this if we need Authentication
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
	user, err := ctrl.userRepository.FindById(1)
	if err != nil {
		http.Error(w, "The user doesn't exist anymore", http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), authenticatedUserKey, user)
	next(r.WithContext(ctx))

	// TODO: Reactivate if we shall use Authentication
	/* token := r.Header.Get("Authorization")

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

	ctx := context.WithValue(r.Context(), authenticatedUserKey, users[0])
	next(r.WithContext(ctx))
	*/
}
