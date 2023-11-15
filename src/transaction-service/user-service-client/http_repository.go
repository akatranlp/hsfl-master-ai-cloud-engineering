package user_service_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/client"
	shared_types "github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/shared-types"
	"net/http"
	"net/url"
)

type HTTPRepository struct {
	userServiceURL *url.URL
	client         client.Client
}

func NewHTTPRepository(userServiceURL *url.URL, client client.Client) *HTTPRepository {
	return &HTTPRepository{userServiceURL, client}
}

func (repo *HTTPRepository) MoveBalance(userId uint64, receivingUserId uint64, amount int64) error {
	host := repo.userServiceURL.String()

	body := &shared_types.MoveBalanceRequest{UserId: userId, ReceivingUserId: receivingUserId, Amount: amount}
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

	if res.StatusCode == http.StatusBadRequest {
		return errors.New("you cannot buy this book")
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("an unknown error")
	}

	var response shared_types.MoveBalanceResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		return errors.New("an unknown error")
	}

	return nil

}
