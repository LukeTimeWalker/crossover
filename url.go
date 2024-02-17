package main

import (
	"fmt"
	"net/url"
)

func LoadUrlParams(u string) (map[string]interface{}, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("could not parse url url=%s error=%s", u, err.Error())
	}

	params, err := url.ParseQuery(ur.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("could not parse url params url=%s error=%s", u, err.Error())
	}

	urlParams := make(map[string]interface{})
	for key, values := range params {
		if len(values) == 1 {
			urlParams[key] = values[0]
		} else {
			urlParams[key] = values
		}
	}

	return urlParams, nil
}
