package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/model"
	chapters_repository "github.com/akatranlp/hsfl-master-ai-cloud-engineering/book-service/chapters/repository"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"golang.org/x/sync/singleflight"
)

type DefaultService struct {
	repository chapters_repository.Repository
	g          *singleflight.Group
}

func NewDefaultService(repository chapters_repository.Repository) *DefaultService {
	g := &singleflight.Group{}
	return &DefaultService{
		repository,
		g,
	}
}

type result struct {
	Chapter         *model.Chapter
	ReceivingUserId *uint64
}

func (s *DefaultService) ValidateChapterId(userId uint64, chapterId uint64, bookId uint64) (*shared_types.ValidateChapterIdResponse, shared_types.Code, error) {
	value, err, _ := s.g.Do(fmt.Sprintf("validate-%d-%d", chapterId, bookId), func() (interface{}, error) {
		chapter, receivingUserId, err := s.repository.ValidateChapterId(chapterId, bookId)
		return &result{
			Chapter:         chapter,
			ReceivingUserId: receivingUserId,
		}, err
	})
	if err != nil {
		log.Println("ERROR [ValidateChapterId - Execute ValidateChapterId]: ", err.Error())
		return nil, shared_types.NotFound, err
	}
	res := value.(*result)
	chapter := res.Chapter
	receivingUserId := res.ReceivingUserId

	if *receivingUserId == userId {
		log.Println("ERROR [ValidateChapterId - receivingUserId == request.UserId]: ", "Author and buyer are the same")
		return nil, shared_types.InvalidArgument, errors.New("author and buyer are the same")
	}

	if chapter.Status != model.Published {
		log.Println("ERROR [ValidateChapterId - chapter.Status != model.Published]: ", "Chapter is not published")
		return nil, shared_types.InvalidArgument, errors.New("chapter is not published")
	}

	return &shared_types.ValidateChapterIdResponse{
			ChapterId:       chapter.ID,
			BookId:          chapter.BookID,
			ReceivingUserId: *receivingUserId,
			Amount:          chapter.Price,
		},
		shared_types.OK,
		nil
}
