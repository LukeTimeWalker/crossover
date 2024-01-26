package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"text/template"

	"github.com/gorilla/mux"
)

type Variables map[string]map[string]interface{}

type PathHandler func(variables Variables) bool

func RegisterRoute(router *mux.Router, route *Route, client *http.Client, envs *map[string]interface{}) {

	variables := make(Variables)

	variables["envs"] = *envs

	pathHandler := make([]PathHandler, len(route.Sender))

	for index, sender := range route.Sender {
		pathHandler[index] = CreatePathHandler(&route.Receiver, &sender, client)
	}

	routeHandler := CreateRouteHandler(route.Receiver, variables, pathHandler, route.Conditions)

	receiverEndpoint := path.Clean(fmt.Sprintf("/%s", route.Receiver.Endpoint))

	router.HandleFunc(receiverEndpoint, routeHandler).Methods(route.Receiver.Method)
}

func CreateRouteHandler(receiver Receiver, variables Variables, pathHandler []PathHandler, conditions []string) func(writer http.ResponseWriter, request *http.Request) {
	logger := GetLogger()

	routeConditionTemplates, err := CreateConditions(conditions)

	if err != nil {
		logger.Error(err.Error())
	}

	return func(writer http.ResponseWriter, request *http.Request) {

		logger.Info("Executing route", "receiverEndpoint", receiver.Endpoint)

		variables["body"] = LoadBody(request)

		if receiver.UrlParams != nil {
			variables["url"] = LoadUrlParams(request.RequestURI, receiver.UrlParams)
		}

		if !EvaluateConditions(routeConditionTemplates, variables) {
			logger.Info("Request was not send due the result of condtions evaluation")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		statusCode := http.StatusOK

		for _, handler := range pathHandler {
			success := handler(variables)
			if !success {
				statusCode = 400
			}
		}

		writer.WriteHeader(statusCode)
	}
}

func CreatePathHandler(receiver *Receiver, sender *Sender, client *http.Client) PathHandler {
	logger := GetLogger()

	var bodyTemplate *template.Template
	var err error

	if sender.Body != "" {
		bodyTemplate, err = CreateTemplate("bodyTemplate", sender.Body)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	var headerTemplates []*template.Template

	for key, value := range sender.Header {
		template, err := CreateTemplate(key, value)
		if err != nil {
			logger.Error(err.Error())
		}
		headerTemplates = append(headerTemplates, template)
	}

	conditionTemplates, err := CreateConditions(sender.Conditions)

	if err != nil {
		logger.Error(err.Error())
	}

	urlTemplate, err := CreateTemplate("urlTemplate", sender.Endpoint)

	if err != nil {
		logger.Error(err.Error())
	}

	pathHandler := func(variables Variables) bool {
		senderBody, err := ExecuteTemplate(bodyTemplate, variables)

		if err != nil {
			logger.Error(err.Error())
			return false
		}

		senderUrl, err := ExecuteTemplate(urlTemplate, variables)

		if err != nil {
			logger.Error(err.Error())
			return false
		}

		// Send

		if !EvaluateConditions(conditionTemplates, variables) {
			logger.Info("Request was not send due the result of condtions evaluation")
			return false
		}

		req, err := http.NewRequest(sender.Method, senderUrl, bytes.NewReader([]byte(senderBody)))

		if err != nil {
			logger.Error("Could not create request", "error", err.Error())
			return false
		}

		if sender.Header != nil {
			for _, tmpl := range headerTemplates {
				header, err := ExecuteTemplate(tmpl, variables)

				if err != nil {
					logger.Error(err.Error())
				}

				req.Header.Add(tmpl.Name(), header)
			}
		}

		logger.Debug("Sending request", "method", sender.Method, "url", senderUrl)

		res, err := client.Do(req)

		if err != nil {
			logger.Error("request failed", "error", err.Error())
			return false
		}

		fmt.Println(res.StatusCode, res.Status)

		if res.StatusCode < http.StatusMultipleChoices && res.StatusCode >= http.StatusOK {
			logger.Info("Request successfull", "status", res.StatusCode)
		} else {
			logger.Error("request unsuccessfull", "status", res.StatusCode)
			return false
		}
		return true
	}

	return pathHandler

}
