package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func LoadBody(request *http.Request) (map[string]interface{}, error) {
	body := make(map[string]interface{})

	err := json.NewDecoder(request.Body).Decode(&body)

	if err != nil {
		return nil, fmt.Errorf("couldn't decode json error=%s", err.Error())
	}

	return body, nil
}
