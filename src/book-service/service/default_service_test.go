package service

import (
	"errors"

	chapters_mocks "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/_mocks/chapters"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"testing"
)

func TestDefaultService(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := chapters_mocks.NewMockRepository(ctrl)
	service := NewDefaultService(repository)

	t.Run("ValidateChapterId", func(t *testing.T) {
		t.Run("should return NotFound if chapter for book is not found", func(t *testing.T) {
			// given
			repository.
				EXPECT().
				ValidateChapterId(uint64(1), uint64(1)).
				Return(nil, nil, errors.New("not found"))

			// when
			response, statusCode, err := service.ValidateChapterId(1, 1, 1)

			// then
			assert.Nil(t, response)
			assert.Equal(t, shared_types.NotFound, statusCode)
			assert.Error(t, err)
		})

		t.Run("should return InvalidArgument if bookAuthor and Buyer are the same", func(t *testing.T) {
			// given
			shouldAuthor := uint64(1)
			shoulChapter := model.Chapter{
				ID:     1,
				BookID: 1,
				Price:  100,
				Status: model.Published,
			}

			repository.
				EXPECT().
				ValidateChapterId(uint64(1), uint64(1)).
				Return(&shoulChapter, &shouldAuthor, nil)

			// when
			response, statusCode, err := service.ValidateChapterId(1, 1, 1)

			// then
			assert.Nil(t, response)
			assert.Equal(t, shared_types.InvalidArgument, statusCode)
			assert.Error(t, err)
		})

		t.Run("should return InvalidArgument if the book is a draft", func(t *testing.T) {
			// given
			shouldAuthor := uint64(2)
			shoulChapter := model.Chapter{
				ID:     1,
				BookID: 1,
				Price:  100,
				Status: model.Draft,
			}

			repository.
				EXPECT().
				ValidateChapterId(uint64(1), uint64(1)).
				Return(&shoulChapter, &shouldAuthor, nil)

			// when
			response, statusCode, err := service.ValidateChapterId(1, 1, 1)

			// then
			assert.Nil(t, response)
			assert.Equal(t, shared_types.InvalidArgument, statusCode)
			assert.Error(t, err)
		})

		t.Run("should return OK if everything works", func(t *testing.T) {
			// given
			shouldAuthor := uint64(2)
			shoulChapter := model.Chapter{
				ID:     1,
				BookID: 1,
				Price:  100,
				Status: model.Published,
			}
			shouldReponse := &shared_types.ValidateChapterIdResponse{
				ChapterId:       shoulChapter.ID,
				BookId:          shoulChapter.BookID,
				ReceivingUserId: shouldAuthor,
				Amount:          shoulChapter.Price,
			}

			repository.
				EXPECT().
				ValidateChapterId(uint64(1), uint64(1)).
				Return(&shoulChapter, &shouldAuthor, nil)

			// when
			response, statusCode, err := service.ValidateChapterId(1, 1, 1)

			// then
			assert.Equal(t, shouldReponse, response)
			assert.Equal(t, shared_types.OK, statusCode)
			assert.NoError(t, err)
		})
	})
}
