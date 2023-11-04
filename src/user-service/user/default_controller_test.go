package user

import (
	"context"
	"encoding/json"
	"errors"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultController(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := mocks.NewMockRepository(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)

	controller := DefaultController{userRepository, hasher, tokenGenerator}

	/* Re-add these tests if we need Authentication again
	t.Run("Authentication-Middleware", func(t *testing.T) {
		t.Run("should return 401 when you don't add a token", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)

			// when
			called := false
			controller.AuthenticationMiddleWare(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 401 when its not a valid token", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			r.Header.Set("Authorization", "tester")

			// when
			called := false
			controller.AuthenticationMiddleWare(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 401 when it's a malformed Bearer token", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			r.Header.Set("Authorization", "Bearer tester")

			// when
			tokenGenerator.EXPECT().VerifyToken("tester").Return(nil, errors.New("token is not valid"))

			called := false
			controller.AuthenticationMiddleWare(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 401 if there is no email claim", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			r.Header.Set("Authorization", "Bearer tester")

			// when
			tokenGenerator.EXPECT().VerifyToken("tester").Return(map[string]interface{}{
				"exp":  12345,
				"user": "tester",
			}, nil)

			called := false
			controller.AuthenticationMiddleWare(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 401 if the user is not found", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			r.Header.Set("Authorization", "Bearer tester")

			// when
			tokenGenerator.EXPECT().VerifyToken("tester").Return(map[string]interface{}{
				"exp":   12345,
				"email": "toni@tester",
			}, nil)

			userRepository.EXPECT().FindByEmail("toni@tester").Return(nil, errors.New("user not found"))

			called := false
			controller.AuthenticationMiddleWare(w, r, func(r *http.Request) {
				called = true
			})

			assert.Equal(t, false, called)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return call next if the token is valid", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			r.Header.Set("Authorization", "Bearer tester")
			user := model.DbUser{Email: "toni@tester"}

			// when
			tokenGenerator.EXPECT().VerifyToken("tester").Return(map[string]interface{}{
				"exp":   12345,
				"email": "toni@tester",
			}, nil)

			userRepository.EXPECT().FindByEmail("toni@tester").Return([]*model.DbUser{&user}, nil)

			called := false
			controller.AuthenticationMiddleWare(w, r, func(req *http.Request) {
				called = true
				r = req
			})

			assert.Equal(t, true, called)
			assert.Equal(t, user, r.Context().Value(authenticatedUserKey))
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
	*/

	t.Run("Login", func(t *testing.T) {
		t.Run("should return 405 METHOD NOT ALLOWED if method is not POST", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/auth/login", nil)

			// when
			controller.Login(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid json`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/auth/login", test)

				// when
				controller.Login(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{"email":"test@test.com"}`),
				strings.NewReader(`{"password":"test"}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/auth/login", test)

				// when
				controller.Login(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if search for user failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return(nil, errors.New("could not query database"))

			// when
			controller.Login(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 401 UNAUTHORIZED if user not found", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"test@test.com","password":"test"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return([]*model.DbUser{}, nil)

			// when
			controller.Login(w, r)

			// then
			assert.Equal(t, "Basic realm=Restricted", w.Header().Get("WWW-Authenticate"))
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 401 UNAUTHORIZED if password is not correct", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"test@test.com","password":"wrong password"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return([]*model.DbUser{{
					Email:    "test@test.com",
					Password: []byte("hashed password"),
				}}, nil)

			hasher.
				EXPECT().
				Validate([]byte("wrong password"), []byte("hashed password")).
				Return(false)

			// when
			controller.Login(w, r)

			// then
			assert.Equal(t, "Basic realm=Restricted", w.Header().Get("WWW-Authenticate"))
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("should return 401 UNAUTHORIZED if password is not correct", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"email":"test@test.com","password":"wrong password"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return([]*model.DbUser{{
					Email:    "test@test.com",
					Password: []byte("hashed password"),
				}}, nil)

			hasher.
				EXPECT().
				Validate([]byte("wrong password"), []byte("hashed password")).
				Return(false)

			// when
			controller.Login(w, r)

			// then
			assert.Equal(t, "Basic realm=Restricted", w.Header().Get("WWW-Authenticate"))
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

	})

	t.Run("Register", func(t *testing.T) {
		t.Run("should return 405 METHOD NOT ALLOWED if method is not POST", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/auth/register", nil)

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid json`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/auth/register", test)

				// when
				controller.Register(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 400 BAD REQUEST if payload is incomplete", func(t *testing.T) {
			tests := []io.Reader{
				strings.NewReader(`{}`),
				strings.NewReader(`{"email":"test@test.com"}`),
				strings.NewReader(`{"password":"test"}`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/api/v1/auth/register", test)

				// when
				controller.Register(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if search for existing user failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(`{"email":"test@test.com","password":"test","profileName":"Toni Tester"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return(nil, errors.New("could not query database"))

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 409 CONFLICT if user already exists", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(`{"email":"test@test.com","password":"test","profileName":"Toni Tester"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return([]*model.DbUser{{}}, nil)

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusConflict, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if hashing password failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(`{"email":"test@test.com","password":"test","profileName":"Toni Tester"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return([]*model.DbUser{}, nil)

			hasher.
				EXPECT().
				Hash([]byte("test")).
				Return(nil, errors.New("could not hash password"))

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if user could be created", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(`{"email":"test@test.com","password":"test","profileName":"Toni Tester"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return([]*model.DbUser{}, nil)

			hasher.
				EXPECT().
				Hash([]byte("test")).
				Return([]byte("hashed password"), nil)

			userRepository.
				EXPECT().
				Create([]*model.DbUser{{
					Email:       "test@test.com",
					Password:    []byte("hashed password"),
					ProfileName: "Toni Tester",
				}}).
				Return(errors.New("could not create user"))

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(`{"email":"test@test.com","password":"test","profileName":"Toni Tester"}`))

			userRepository.
				EXPECT().
				FindByEmail("test@test.com").
				Return([]*model.DbUser{}, nil)

			hasher.
				EXPECT().
				Hash([]byte("test")).
				Return([]byte("hashed password"), nil)

			userRepository.
				EXPECT().
				Create([]*model.DbUser{{
					Email:       "test@test.com",
					Password:    []byte("hashed password"),
					ProfileName: "Toni Tester",
				}}).
				Return(nil)

			// when
			controller.Register(w, r)

			// then
			assert.Equal(t, http.StatusCreated, w.Code)
		})
	})

	t.Run("GetUsers", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users", nil)

			userRepository.
				EXPECT().
				FindAll().
				Return(nil, errors.New("query failed")).
				Times(1)

			// when
			controller.GetUsers(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return all products", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users", nil)

			userRepository.
				EXPECT().
				FindAll().
				Return([]*model.DbUser{{ID: 999}}, nil).
				Times(1)

			// when
			controller.GetUsers(w, r)

			// then
			res := w.Result()
			var response []model.UserDTO
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Len(t, response, 1)
			assert.Equal(t, uint64(999), response[0].ID)
		})
	})

	t.Run("GetUser", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "userid", "1"))

			userRepository.
				EXPECT().
				FindById(uint64(1)).
				Return(nil, errors.New("database error"))

			// when
			controller.GetUser(w, r)

			// then
			assert.Equal(t, http.StatusNotFound, w.Code)
		})

		t.Run("should return 200 OK and user", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/1", nil)
			r = r.WithContext(context.WithValue(r.Context(), "userid", "1"))
			id := uint64(1)

			userRepository.
				EXPECT().
				FindById(id).
				Return(&model.DbUser{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("hash"),
					ProfileName: "Toni Tester",
					Balance:     0,
				}, nil)

			// when
			controller.GetUser(w, r)

			// then
			res := w.Result()
			var response model.UserDTO
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, id, response.ID)
			assert.Equal(t, "test@test.com", response.Email)
		})
	})

	//===================================================================================================

	t.Run("GetMe", func(t *testing.T) {
		t.Run("Should return 200 when its called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)
			dbUser := &model.DbUser{
				ID:          1,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			}
			r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

			// when
			controller.GetMe(w, r)

			// then
			res := w.Result()
			var response model.UserDTO
			err := json.NewDecoder(res.Body).Decode(&response)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, uint64(1), response.ID)
			assert.Equal(t, "test@test.com", response.Email)

		})
	})

	t.Run("PatchMe", func(t *testing.T) {
		t.Run("should return 400 BAD REQUEST if payload is not json", func(t *testing.T) {
			tests := []io.Reader{
				nil,
				strings.NewReader(`{"invalid`),
			}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest("PUT", "/api/v1/users/me", test)
				dbUser := &model.DbUser{
					ID:          1,
					Email:       "test@test.com",
					Password:    []byte("hash"),
					ProfileName: "Toni Tester",
					Balance:     0,
				}
				r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

				// when
				controller.PatchMe(w, r)

				// then
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if query failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/me",
				strings.NewReader(`{"profileName":"Tino Taster"}`))
			dbUser := &model.DbUser{
				ID:          1,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			}
			r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

			userRepository.
				EXPECT().
				Update(uint64(1), gomock.Any()).
				Do(func(_ uint64, user *model.DbUserPatch) {
					assert.Equal(t, "Tino Taster", *user.ProfileName)
				}).
				Return(errors.New("database error"))

			// when
			controller.PatchMe(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 and update profile Name", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/me",
				strings.NewReader(`{"profileName":"Tino Taster"}`))
			dbUser := &model.DbUser{
				ID:          1,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			}
			r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

			userRepository.
				EXPECT().
				Update(uint64(1), gomock.Any()).
				Do(func(_ uint64, user *model.DbUserPatch) {
					assert.Equal(t, "Tino Taster", *user.ProfileName)
				}).
				Return(nil)

			// when
			controller.PatchMe(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should return 500 if hashing failed", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/me",
				strings.NewReader(`{"profileName":"Tino Taster","password":"test2"}`))
			dbUser := &model.DbUser{
				ID:          1,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			}
			r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

			hasher.
				EXPECT().
				Hash([]byte("test2")).
				Return(nil, errors.New("could not hash password"))

			// when
			controller.PatchMe(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/users/me",
				strings.NewReader(`{"profileName":"Tino Taster","password":"test2"}`))
			dbUser := &model.DbUser{
				ID:          1,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			}
			r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

			hasher.
				EXPECT().
				Hash([]byte("test2")).
				Return([]byte("hashed password"), nil)

			userRepository.
				EXPECT().
				Update(uint64(1), gomock.Any()).
				Do(func(_ uint64, user *model.DbUserPatch) {
					assert.Equal(t, "Tino Taster", *user.ProfileName)
					assert.Equal(t, []byte("hashed password"), *user.Password)
				}).
				Return(nil)

			// when
			controller.PatchMe(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("DeleteMe", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if query fails", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/users/me", nil)
			dbUser := &model.DbUser{
				ID:          1,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			}
			r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

			userRepository.
				EXPECT().
				Delete([]*model.DbUser{dbUser}).
				Return(errors.New("database error"))

			// when
			controller.DeleteMe(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/users/me", nil)
			dbUser := &model.DbUser{
				ID:          1,
				Email:       "test@test.com",
				Password:    []byte("hash"),
				ProfileName: "Toni Tester",
				Balance:     0,
			}
			r = r.WithContext(context.WithValue(r.Context(), authenticatedUserKey, dbUser))

			userRepository.
				EXPECT().
				Delete([]*model.DbUser{dbUser}).
				Return(nil)

			// when
			controller.DeleteMe(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
