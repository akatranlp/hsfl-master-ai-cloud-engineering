package auth_middleware

import (
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib/router"
	"net/http"
)

type Controller interface {
	AuthenticationMiddleware(http.ResponseWriter, *http.Request, router.Next)
}
