package auth_middleware

import (
	"net/http"

	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
)

type Controller interface {
	AuthenticationMiddleware(http.ResponseWriter, *http.Request, router.Next)
}
