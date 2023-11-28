package balancer

import (
	"net/http"
	"net/url"
)

func GetHealth(client *http.Client, host *url.URL) bool {
	resp, err := client.Get(host.Scheme + "://" + host.Host + "/health")
	return err == nil && resp.StatusCode == http.StatusOK
}
