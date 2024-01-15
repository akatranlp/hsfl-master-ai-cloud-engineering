package controller

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/test-data-service/repository"
)

type DefaultController struct {
	repository repository.Repository
}

func NewDefaultController(repository repository.Repository) *DefaultController {
	return &DefaultController{repository}
}

func (c *DefaultController) ResetDatabase(w http.ResponseWriter, r *http.Request) {
	err := c.repository.ResetDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
