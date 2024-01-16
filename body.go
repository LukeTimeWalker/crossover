package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func CompareBodies(configBody map[string]interface{}, body map[string]interface{}) {
	for key, value := range configBody {
		if reflect.ValueOf(value).Kind() == reflect.Map {
			CompareBodies(configBody[key].(map[string]interface{}), body[key].(map[string]interface{}))
		} else if reflect.ValueOf(value).Kind() == reflect.Array {
			for _, a := range body[key].([]interface{}) {
				CompareBodies(configBody[key].(map[string]interface{}), a.(map[string]interface{}))
			}
		} else {
			if _, ok := body[key]; ok {
				configBody[key] = body[key]
			}
		}
	}
}

func GenerateBodyStruct(input map[string]interface{}) any {
	fields := []reflect.StructField{}

	for key, value := range input {
		field := reflect.StructField{
			Name: cases.Title(language.AmericanEnglish).String(key),
			Type: reflect.TypeOf(value),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, key)),
		}

		fields = append(fields, field)
	}

	newStruct := reflect.StructOf(fields)

	return reflect.New(newStruct).Interface()

}
