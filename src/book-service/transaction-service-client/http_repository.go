package transaction_service_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/client"
	"net/http"
	"net/url"

	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
)

type HTTPRepository struct {
	transactionServiceURL *url.URL
	client                client.Client
}

func NewHTTPRepository(transactionServiceURL *url.URL, client client.Client) *HTTPRepository {
	return &HTTPRepository{transactionServiceURL, client}
}

func (repo *HTTPRepository) CheckChapterBought(userId uint64, chapterId uint64) error {
	host := repo.transactionServiceURL.String()

	body := &shared_types.CheckChapterBoughtRequest{UserID: userId, ChapterID: chapterId}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", host, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	res, err := repo.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return errors.New("you haven't bought this book")
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("an unknown error")
	}

	var response shared_types.CheckChapterBoughtResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		return errors.New("an unknown error")
	}

	return nil
}
