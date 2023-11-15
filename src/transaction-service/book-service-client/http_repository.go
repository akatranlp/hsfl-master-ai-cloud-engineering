package book_service_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/client"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"net/http"
	"net/url"
)

type validateChapterIdRequest struct {
	UserId    uint64 `json:"userId"`
	ChapterId uint64 `json:"chapterId"`
}

type HTTPRepository struct {
	bookServiceURL *url.URL
	client         client.Client
}

func NewHTTPRepository(bookServiceURL *url.URL, client client.Client) *HTTPRepository {
	return &HTTPRepository{bookServiceURL, client}
}

func (repo *HTTPRepository) ValidateChapterId(userId uint64, chapterId uint64) (*shared_types.ValidateChapterIdResponse, error) {
	host := repo.bookServiceURL.String()

	body := &validateChapterIdRequest{userId, chapterId}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", host, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	res, err := repo.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusBadRequest {
		return nil, errors.New("you cannot buy this book")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("an unknown error")
	}

	var response shared_types.ValidateChapterIdResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil

}
