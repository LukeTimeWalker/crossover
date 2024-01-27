package main

import (
	"bytes"
	"net/http"
	"text/template"
)

type Variables map[string]map[string]interface{}

type PathHandler func(variables Variables) bool

func CreateRouteHandler(route *Route, envs *map[string]interface{}, client *http.Client) (func(writer http.ResponseWriter, request *http.Request), error) {

	variables := make(Variables)

	variables["envs"] = *envs

	pathHandler := make([]PathHandler, len(route.Sender))

	for index, sender := range route.Sender {
		handler, err := CreatePathHandler(&route.Receiver, &sender, client)
		if err != nil {
			return nil, err
		}
		pathHandler[index] = handler
	}

	routeConditionTemplates, err := CreateConditions(route.Conditions)

	if err != nil {
		return nil, err
	}

	routeHandler := func(writer http.ResponseWriter, request *http.Request) {

		logger := GetLogger()

		logger.Info("Executing route", "receiverEndpoint", route.Receiver.Endpoint)

		variables["body"] = LoadBody(request)

		if route.Receiver.UrlParams != nil {
			variables["url"] = LoadUrlParams(request.RequestURI, route.Receiver.UrlParams)
		}

		conditionResult, err := EvaluateConditions(routeConditionTemplates, variables)

		if err != nil {
			logger.Error(err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if !conditionResult {
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

	return routeHandler, nil
}

func CreatePathHandler(receiver *Receiver, sender *Sender, client *http.Client) (PathHandler, error) {

	var bodyTemplate *template.Template
	var err error

	if sender.Body != "" {
		bodyTemplate, err = CreateTemplate("bodyTemplate", sender.Body)
		if err != nil {
			return nil, err
		}
	}

	var headerTemplates []*template.Template

	for key, value := range sender.Header {
		template, err := CreateTemplate(key, value)
		if err != nil {
			return nil, err
		}
		headerTemplates = append(headerTemplates, template)
	}

	conditionTemplates, err := CreateConditions(sender.Conditions)

	if err != nil {
		return nil, err
	}

	urlTemplate, err := CreateTemplate("urlTemplate", sender.Endpoint)

	if err != nil {
		return nil, err
	}

	pathHandler := func(variables Variables) bool {

		logger := GetLogger()

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

		conditionResult, err := EvaluateConditions(conditionTemplates, variables)

		if err != nil {
			logger.Error(err.Error())
			return false
		}
		if !conditionResult {
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

		if res.StatusCode < http.StatusMultipleChoices && res.StatusCode >= http.StatusOK {
			logger.Info("Request successfull", "status", res.StatusCode)
		} else {
			logger.Error("request unsuccessfull", "status", res.StatusCode)
			return false
		}
		return true
	}

	return pathHandler, nil

}
