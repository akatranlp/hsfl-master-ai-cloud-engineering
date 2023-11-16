package health

import "net/http"

type DefaultController struct{}

func NewDefaultController() *DefaultController {
	return &DefaultController{}
}

func (ctrl *DefaultController) ProvideHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
