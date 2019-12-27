package analyze

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func get(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to get tickstore response. %+v", err)
	}
	return resp, nil
}

func getWithAuth(url, username, password string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to get tickstore response with basic auth. %+v", err)
	}
	return resp, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
