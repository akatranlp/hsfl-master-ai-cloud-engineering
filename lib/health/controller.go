package health

import "net/http"

type Controller interface {
	ProvideHealth(http.ResponseWriter, *http.Request)
}
