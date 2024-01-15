package service

import (
	"errors"

	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/_mocks"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/transaction-service/model"
	"github.com/stretchr/testify/assert"

	"testing"

	"go.uber.org/mock/gomock"
)

func TestDefaultService(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mocks.NewMockRepository(ctrl)
	service := NewDefaultService(repository)

	t.Run("CheckChapterBought", func(t *testing.T) {
		t.Run("should return false if transaction is not found", func(t *testing.T) {
			// given
			repository.
				EXPECT().
				FindForUserIdAndChapterId(uint64(1), uint64(1), uint64(1)).
				Return(nil, errors.New("not found"))

			// when
			success, statusCode, err := service.CheckChapterBought(1, 1, 1)

			// given
			assert.Error(t, err)
			assert.Equal(t, shared_types.NotFound, statusCode)
			assert.False(t, success)
		})

		t.Run("should return false if transaction is nil", func(t *testing.T) {
			// given
			repository.
				EXPECT().
				FindForUserIdAndChapterId(uint64(1), uint64(1), uint64(1)).
				Return(nil, nil)

			// when
			success, statusCode, err := service.CheckChapterBought(1, 1, 1)

			// given
			assert.NoError(t, err)
			assert.Equal(t, shared_types.OK, statusCode)
			assert.False(t, success)
		})

		t.Run("should return true if transaction is found", func(t *testing.T) {
			// given
			transaction := &model.Transaction{
				ID:              1,
				BookID:          1,
				ChapterID:       1,
				PayingUserID:    2,
				ReceivingUserID: 1,
				Amount:          100,
			}

			repository.
				EXPECT().
				FindForUserIdAndChapterId(uint64(1), uint64(1), uint64(1)).
				Return(transaction, nil)

			// when
			success, statusCode, err := service.CheckChapterBought(1, 1, 1)

			// given
			assert.NoError(t, err)
			assert.Equal(t, shared_types.OK, statusCode)
			assert.True(t, success)
		})
	})
}
