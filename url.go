package main

import (
	"fmt"
	"net/url"
	"strconv"
)

func LoadUrlParams(u string, urlParams map[string]string) (map[string]interface{}, error) {
	ur, err := url.Parse(u)

	if err != nil {
		return nil, fmt.Errorf("could not parse url url=%s error=%s", u, err.Error())
	}

	params, err := url.ParseQuery(ur.RawQuery)

	if err != nil {
		return nil, fmt.Errorf("could not parse url params url=%s error=%s", u, err.Error())
	}

	output := make(map[string]interface{})

	for key, value := range urlParams {
		param, err := parseUrlParam(key, params.Get(key), value)

		if err != nil {
			return nil, err
		}

		output[key] = param
	}

	return output, nil
}

func parseUrlParam(key string, value string, t string) (interface{}, error) {
	var returnValue interface{}

	var err error

	switch t {
	case "string":
		returnValue = value
	case "int":
		returnValue, err = strconv.ParseInt(value, 10, 0)
	case "boolean":
		returnValue, err = strconv.ParseBool(value)
	case "float":
		returnValue, err = strconv.ParseFloat(value, 64)
	default:
		return nil, fmt.Errorf("unknown type for url param type=%s", t)
	}

	if err != nil {
		return nil, fmt.Errorf("could not parse url param param=%s value=%s expectedType=%s", key, value, t)
	}

	return returnValue, nil
}
