package main

import (
	"net/url"
	"strconv"
)

func LoadUrlParams(u string, urlParams map[string]string) map[string]interface{} {
	logger := GetLogger()

	ur, err := url.Parse(u)

	if err != nil {
		logger.Error("Could not parse url", "url", u, "error", err)
	}

	params, err := url.ParseQuery(ur.RawQuery)

	if err != nil {
		logger.Error("Could not parse url params", "url", u, "error", err)
	}

	output := make(map[string]interface{})

	for key, value := range urlParams {
		output[key] = parseUrlParam(key, params.Get(key), value)
	}

	return output
}

func parseUrlParam(key string, value string, t string) interface{} {

	logger := GetLogger()

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
		logger.Error("Unknown type for url param", "type", t)
		returnValue = nil
	}
	if err != nil {
		logger.Error("Could not parse url param", "param", key, "value", value, "expectedType", t)
	}
	return returnValue
}
