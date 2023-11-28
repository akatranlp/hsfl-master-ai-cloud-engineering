package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	health_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/health/_mocks"
	libRouter "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)

	userController := mocks.NewMockController(ctrl)
	healthController := health_mocks.NewMockController(ctrl)
	router := New(userController, healthController)

	t.Run("/api/v1/users", func(t *testing.T) {
		t.Run("GetUsers should not be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users", nil)

			userController.
				EXPECT().
				AuthenticationMiddleWare(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					w.WriteHeader(http.StatusUnauthorized)
				}).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)

		})

		t.Run("GetUsers should be called", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users", nil)

			userController.
				EXPECT().
				AuthenticationMiddleWare(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).
				Times(1)

			userController.
				EXPECT().
				GetUsers(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

		})
	})

	t.Run("/api/v1/users", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET", func(t *testing.T) {
			tests := []string{"HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/users", nil)

				userController.
					EXPECT().
					AuthenticationMiddleWare(w, r, gomock.Any()).
					Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
						next(r)
					}).
					Times(1)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users", nil)

			userController.
				EXPECT().
				AuthenticationMiddleWare(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).
				Times(1)

			userController.
				EXPECT().
				GetUsers(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/login", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {

			tests := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/login", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/login", nil)

			userController.
				EXPECT().
				Login(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/register", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/register", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/register", nil)

			userController.
				EXPECT().
				Register(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/refresh-token", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/refresh-token", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/refresh-token", nil)

			userController.
				EXPECT().
				RefreshToken(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/logout", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/logout", nil)

				userController.
					EXPECT().
					AuthenticationMiddleWare(w, r, gomock.Any()).
					Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
						next(r)
					}).
					Times(1)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/logout", nil)

			userController.
				EXPECT().
				AuthenticationMiddleWare(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).
				Times(1)

			userController.
				EXPECT().
				Logout(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/api/v1/users/me", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not GET, DELETE or PATCH", func(t *testing.T) {
			tests := []string{"HEAD", "POST", "CONNECT", "OPTIONS", "TRACE", "PUT"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/api/v1/users/me", nil)

				userController.
					EXPECT().
					AuthenticationMiddleWare(w, r, gomock.Any()).
					Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
						next(r)
					}).
					Times(1)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call GET handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/users/me", nil)

			userController.
				EXPECT().
				AuthenticationMiddleWare(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).
				Times(1)

			userController.
				EXPECT().
				GetMe(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call PATCH handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PATCH", "/api/v1/users/me", nil)

			userController.
				EXPECT().
				AuthenticationMiddleWare(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).
				Times(1)

			userController.
				EXPECT().
				PatchMe(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})

		t.Run("should call Delete handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/users/me", nil)

			userController.
				EXPECT().
				AuthenticationMiddleWare(w, r, gomock.Any()).
				Do(func(w http.ResponseWriter, r *http.Request, next libRouter.Next) {
					next(r)
				}).
				Times(1)

			userController.
				EXPECT().
				DeleteMe(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/validate-token", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"HEAD", "GET", "CONNECT", "OPTIONS", "TRACE", "PATCH", "DELETE", "PUT"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/validate-token", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/validate-token", nil)

			userController.
				EXPECT().
				ValidateToken(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})

	t.Run("/move-user-amount", func(t *testing.T) {
		t.Run("should return 404 NOT FOUND if method is not POST", func(t *testing.T) {
			tests := []string{"HEAD", "GET", "CONNECT", "OPTIONS", "TRACE", "PATCH", "DELETE", "PUT"}

			for _, test := range tests {
				// given
				w := httptest.NewRecorder()
				r := httptest.NewRequest(test, "/move-user-amount", nil)

				// when
				router.ServeHTTP(w, r)

				// then
				assert.Equal(t, http.StatusNotFound, w.Code)
			}
		})

		t.Run("should call POST handler", func(t *testing.T) {
			// given
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/move-user-amount", nil)

			userController.
				EXPECT().
				MoveUserAmount(w, r).
				Times(1)

			// when
			router.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Code)
		})
	})
}
