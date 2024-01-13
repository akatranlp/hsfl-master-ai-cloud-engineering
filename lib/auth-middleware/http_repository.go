package auth_middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/client"
)

type HTTPRepository struct {
	authServiceURL *url.URL
	client         client.Client
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}

type VerifyTokenResponse struct {
	Success bool   `json:"success"`
	UserId  uint64 `json:"userId"`
}

func NewHTTPRepository(authServiceURL *url.URL, client client.Client) *HTTPRepository {
	return &HTTPRepository{authServiceURL, client}
}

func (repo *HTTPRepository) VerifyToken(token string) (uint64, error) {
	host := repo.authServiceURL.String()
	tokenBody := &VerifyTokenRequest{token}

	reqBody, err := json.Marshal(tokenBody)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", host, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, err
	}

	res, err := repo.client.Do(req)
	if err != nil {
		return 0, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return 0, errors.New("you are not authorized to do this")
	}

	if res.StatusCode != http.StatusOK {
		return 0, errors.New("an unknown error")
	}

	var response VerifyTokenResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, err
	}

	if !response.Success {
		return 0, errors.New("an unknown error")
	}

	return response.UserId, nil
}
