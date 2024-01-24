package main

import (
	"encoding/json"
	"net/http"
)

func LoadBody(request *http.Request) map[string]interface{} {

	logger := GetLogger()

	body := make(map[string]interface{})

	err := json.NewDecoder(request.Body).Decode(&body)

	if err != nil {
		logger.Error("Couldn't decode json", "error", err.Error())
	}

	return body
}
