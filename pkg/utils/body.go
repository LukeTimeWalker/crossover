package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func LoadBody(requestBody io.ReadCloser) (map[string]interface{}, error) {
	if requestBody == http.NoBody {
		return nil, nil
	}

	body := make(map[string]interface{})

	bytes, err := io.ReadAll(requestBody)

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
