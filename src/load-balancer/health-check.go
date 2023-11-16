package main

import "net/http"

func IsServiceHealthy() (bool, error) {
	resp, err := http.Get("")
	if err != nil {
		return false, nil
	} else if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, err
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	health, err := IsServiceHealthy()
	if err == nil && health {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusServiceUnavailable)
}
