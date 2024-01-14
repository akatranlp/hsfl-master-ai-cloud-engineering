package service

import (
	"errors"
	"testing"

	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/_mocks"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/user-service/user/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDefaultService(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mocks.NewMockRepository(ctrl)
	tokenGenerator := mocks.NewMockTokenGenerator(ctrl)
	service := NewDefaultService(repository, tokenGenerator, true)

	t.Run("Auth Diactivated", func(t *testing.T) {
		service := NewDefaultService(repository, tokenGenerator, false)

		t.Run("ValidateToken", func(t *testing.T) {
			t.Run("return NotFound if error accured", func(t *testing.T) {
				// given
				repository.
					EXPECT().
					FindById(uint64(1)).
					Return(nil, errors.New("Not found"))

				// when
				user, statusCode, err := service.ValidateToken("")

				// then
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, shared_types.NotFound, statusCode)
			})

			t.Run("return Ok if user was found", func(t *testing.T) {
				// given
				shouldUser := &model.DbUser{
					ID: 1,
				}
				repository.
					EXPECT().
					FindById(uint64(1)).
					Return(shouldUser, nil)

				// when
				user, statusCode, err := service.ValidateToken("")

				// then
				assert.NoError(t, err)
				assert.Equal(t, shouldUser, user)
				assert.Equal(t, shared_types.OK, statusCode)
			})
		})
	})

	t.Run("Auth Activated", func(t *testing.T) {

		t.Run("ValidateToken", func(t *testing.T) {
			t.Run("return Unauthenticated if token is invalid", func(t *testing.T) {
				// given
				tokenGenerator.
					EXPECT().
					VerifyToken("token").
					Return(nil, errors.New("Unauthenticated"))

				// when
				user, statusCode, err := service.ValidateToken("token")

				// then
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, shared_types.Unauthenticated, statusCode)
			})

			t.Run("return Unauthenticated if email claim is missing", func(t *testing.T) {
				// given
				claims := map[string]interface{}{}

				tokenGenerator.
					EXPECT().
					VerifyToken("token").
					Return(claims, nil)

				// when
				user, statusCode, err := service.ValidateToken("token")

				// then
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, shared_types.Unauthenticated, statusCode)
			})

			t.Run("return Unauthenticated if tokenVersion is missing", func(t *testing.T) {
				// given
				claims := map[string]interface{}{
					"email": "test@test.com",
				}

				tokenGenerator.
					EXPECT().
					VerifyToken("token").
					Return(claims, nil)

				// when
				user, statusCode, err := service.ValidateToken("token")

				// then
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, shared_types.Unauthenticated, statusCode)
			})

			t.Run("return Internal if db-Error accured", func(t *testing.T) {
				// given
				claims := map[string]interface{}{
					"email":         "test@test.com",
					"token_version": float64(0),
				}

				tokenGenerator.
					EXPECT().
					VerifyToken("token").
					Return(claims, nil)

				repository.
					EXPECT().
					FindByEmail("test@test.com").
					Return(nil, errors.New("internal error"))

				// when
				user, statusCode, err := service.ValidateToken("token")

				// then
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, shared_types.Internal, statusCode)
			})

			t.Run("return Unauthenticated if no users were found", func(t *testing.T) {
				// given
				claims := map[string]interface{}{
					"email":         "test@test.com",
					"token_version": float64(0),
				}

				users := []*model.DbUser{}

				tokenGenerator.
					EXPECT().
					VerifyToken("token").
					Return(claims, nil)

				repository.
					EXPECT().
					FindByEmail("test@test.com").
					Return(users, nil)

				// when
				user, statusCode, err := service.ValidateToken("token")

				// then
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, shared_types.Unauthenticated, statusCode)
			})

			t.Run("return Unauthenticated if token Version of the user and token are different", func(t *testing.T) {
				// given
				claims := map[string]interface{}{
					"email":         "test@test.com",
					"token_version": float64(0),
				}

				users := []*model.DbUser{
					{
						ID:           1,
						TokenVersion: 1,
					},
				}

				tokenGenerator.
					EXPECT().
					VerifyToken("token").
					Return(claims, nil)

				repository.
					EXPECT().
					FindByEmail("test@test.com").
					Return(users, nil)

				// when
				user, statusCode, err := service.ValidateToken("token")

				// then
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.Equal(t, shared_types.Unauthenticated, statusCode)
			})

			t.Run("return OK if no error accured", func(t *testing.T) {
				// given
				claims := map[string]interface{}{
					"email":         "test@test.com",
					"token_version": float64(0),
				}

				users := []*model.DbUser{
					{
						ID:           1,
						TokenVersion: 0,
					},
				}

				tokenGenerator.
					EXPECT().
					VerifyToken("token").
					Return(claims, nil)

				repository.
					EXPECT().
					FindByEmail("test@test.com").
					Return(users, nil)

				// when
				user, statusCode, err := service.ValidateToken("token")

				// then
				assert.NoError(t, err)
				assert.Equal(t, users[0], user)
				assert.Equal(t, shared_types.OK, statusCode)
			})
		})
	})

	t.Run("MoveUserAmount", func(t *testing.T) {
		t.Run("return NotFound if payingUser was not Found", func(t *testing.T) {
			// given
			payingUserId := uint64(1)
			receivingUserId := uint64(2)

			repository.
				EXPECT().
				FindById(payingUserId).
				Return(nil, errors.New("Not found"))

			// when
			statusCode, err := service.MoveUserAmount(payingUserId, receivingUserId, 100)

			// then
			assert.Error(t, err)
			assert.Equal(t, shared_types.NotFound, statusCode)
		})

		t.Run("return NotFound if receivingUser was not Found", func(t *testing.T) {
			// given
			payingUserId := uint64(1)
			payingUser := &model.DbUser{
				ID:      payingUserId,
				Balance: 1000,
			}

			receivingUserId := uint64(2)

			repository.
				EXPECT().
				FindById(payingUserId).
				Return(payingUser, nil)

			repository.
				EXPECT().
				FindById(receivingUserId).
				Return(nil, errors.New("Not found"))

			// when
			statusCode, err := service.MoveUserAmount(payingUserId, receivingUserId, 100)

			// then
			assert.Error(t, err)
			assert.Equal(t, shared_types.NotFound, statusCode)
		})

		t.Run("return Internal if payingUser update error accured", func(t *testing.T) {
			// given
			payingUserId := uint64(1)
			payingUser := &model.DbUser{
				ID:      payingUserId,
				Balance: 1000,
			}
			payingUserNewBalance := payingUser.Balance - 100

			receivingUserId := uint64(2)
			receivingUser := &model.DbUser{
				ID:      receivingUserId,
				Balance: 2000,
			}

			repository.
				EXPECT().
				FindById(payingUserId).
				Return(payingUser, nil)

			repository.
				EXPECT().
				FindById(receivingUserId).
				Return(receivingUser, nil)

			repository.
				EXPECT().
				Update(payingUserId, &model.DbUserPatch{Balance: &payingUserNewBalance}).
				Return(errors.New("Internal error"))

			// when
			statusCode, err := service.MoveUserAmount(payingUserId, receivingUserId, 100)

			// then
			assert.Error(t, err)
			assert.Equal(t, shared_types.Internal, statusCode)
		})

		t.Run("return Internal if receivingUser update error accured", func(t *testing.T) {
			// given
			payingUserId := uint64(1)
			payingUser := &model.DbUser{
				ID:      payingUserId,
				Balance: 1000,
			}
			payingUserNewBalance := payingUser.Balance - 100

			receivingUserId := uint64(2)
			receivingUser := &model.DbUser{
				ID:      receivingUserId,
				Balance: 2000,
			}
			receivingUserNewBalance := receivingUser.Balance + 100

			repository.
				EXPECT().
				FindById(payingUserId).
				Return(payingUser, nil)

			repository.
				EXPECT().
				FindById(receivingUserId).
				Return(receivingUser, nil)

			repository.
				EXPECT().
				Update(payingUserId, &model.DbUserPatch{Balance: &payingUserNewBalance}).
				Return(nil)

			repository.
				EXPECT().
				Update(receivingUserId, &model.DbUserPatch{Balance: &receivingUserNewBalance}).
				Return(errors.New("Internal error"))

			// when
			statusCode, err := service.MoveUserAmount(payingUserId, receivingUserId, 100)

			// then
			assert.Error(t, err)
			assert.Equal(t, shared_types.Internal, statusCode)
		})

		t.Run("return Ok if no error accured", func(t *testing.T) {
			// given
			payingUserId := uint64(1)
			payingUser := &model.DbUser{
				ID:      payingUserId,
				Balance: 1000,
			}
			payingUserNewBalance := payingUser.Balance - 100

			receivingUserId := uint64(2)
			receivingUser := &model.DbUser{
				ID:      receivingUserId,
				Balance: 2000,
			}
			receivingUserNewBalance := receivingUser.Balance + 100

			repository.
				EXPECT().
				FindById(payingUserId).
				Return(payingUser, nil)

			repository.
				EXPECT().
				FindById(receivingUserId).
				Return(receivingUser, nil)

			repository.
				EXPECT().
				Update(payingUserId, &model.DbUserPatch{Balance: &payingUserNewBalance}).
				Return(nil)

			repository.
				EXPECT().
				Update(receivingUserId, &model.DbUserPatch{Balance: &receivingUserNewBalance}).
				Return(nil)

			// when
			statusCode, err := service.MoveUserAmount(payingUserId, receivingUserId, 100)

			// then
			assert.NoError(t, err)
			assert.Equal(t, shared_types.OK, statusCode)
		})
	})
}
