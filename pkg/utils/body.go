package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func LoadBody(request *http.Request) (map[string]interface{}, error) {
	body := make(map[string]interface{})

	bytes, err := io.ReadAll(request.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read body error=%s", err.Error())
	}

	str := string(bytes)

	str = strings.ReplaceAll(str, "\\n", "\\\\n")

	err = json.Unmarshal([]byte(str), &body)

	if err != nil {
		return nil, fmt.Errorf("failed to decode json error=%s", err.Error())
	}

	return body, nil
}
