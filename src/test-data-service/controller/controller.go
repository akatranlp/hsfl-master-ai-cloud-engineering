package controller

import "net/http"

type Controller interface {
	ResetDatabase(http.ResponseWriter, *http.Request)
}
