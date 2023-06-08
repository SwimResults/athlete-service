package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Post(url string, path string, data interface{}) (*http.Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest("POST", url+path, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}
