package analyze

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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

//postWithAuth send analyser result to result store api to store in mongodb
func postWithAuth(url, username, password string, res *Result) (*http.Response, error) {
	client := &http.Client{}
	reqBody, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("error marshalling ohlc analysis result. %+v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to store result in resultStore with basic auth. %+v", err)
	}
	return resp, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
