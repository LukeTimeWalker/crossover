package main

import (
	"bytes"
	"net/http"
	"text/template"
)

type Variables map[string]map[string]interface{}

type PathHandler func(variables Variables) bool

type RouteHandler func(writer http.ResponseWriter, request *http.Request)

func CreateRouteHandler(route *RouteConfig, envs *map[string]interface{}) (RouteHandler, error) {

	pathHandlers := make([]PathHandler, len(route.Sender))

	for index, sender := range route.Sender {
		handler, err := CreatePathHandler(&route.Receiver, &sender)
		if err != nil {
			return nil, err
		}
		pathHandlers[index] = handler
	}

	routeConditionTemplates, err := CreateConditions(route.Conditions)

	if err != nil {
		return nil, err
	}

	routeHandler := func(writer http.ResponseWriter, request *http.Request) {

		logger := GetLogger()

		logger.Info("Executing route", "receiverEndpoint", route.Receiver.Endpoint)

		body := make(map[string]interface{})

		if request.Body != http.NoBody {
			body, err = LoadBody(request)

			if err != nil {
				logger.Error(err.Error())
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		urlParams, err := LoadUrlParams(request.RequestURI)

		if err != nil {
			logger.Error(err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		variables := Variables{
			"body": body,
			"url":  urlParams,
			"env":  *envs,
		}

		conditionResult, err := EvaluateConditions(routeConditionTemplates, variables)

		if err != nil {
			logger.Error(err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if !conditionResult {
			logger.Info("Request was not sent due to the result of condtions evaluation")
			writer.WriteHeader(http.StatusOK)
			return
		}

		statusCode := http.StatusOK

		for _, handler := range pathHandlers {
			success := handler(variables)
			if !success {
				statusCode = http.StatusBadRequest
			}
		}

		writer.WriteHeader(statusCode)
	}

	return routeHandler, nil
}

func CreatePathHandler(receiverConfig *ReceiverConfig, senderConfig *SenderConfig) (PathHandler, error) {
	bodyTemplate, err := CreateTemplate("bodyTemplate", senderConfig.Body)
	if err != nil {
		return nil, err
	}

	headerTemplates := make(map[string]*template.Template)
	for key, value := range senderConfig.Header {
		template, err := CreateTemplate(key, value)
		if err != nil {
			return nil, err
		}
		headerTemplates[key] = template
	}

	conditionTemplates, err := CreateConditions(senderConfig.Conditions)
	if err != nil {
		return nil, err
	}

	urlTemplate, err := CreateTemplate("urlTemplate", senderConfig.Endpoint)
	if err != nil {
		return nil, err
	}

	client, err := GetClient(senderConfig.Client)
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
			logger.Info("Request was not sent due to the result of conditions evaluation")
			return true
		}

		req, err := http.NewRequest(senderConfig.Method, senderUrl, bytes.NewReader([]byte(senderBody)))
		if err != nil {
			logger.Error("Could not create request", "error", err.Error())
			return false
		}

		for key, tmpl := range headerTemplates {
			header, err := ExecuteTemplate(tmpl, variables)
			if err != nil {
				logger.Error(err.Error())
			}
			req.Header.Add(key, header)
		}

		logger.Debug("Sending request", "method", senderConfig.Method, "url", senderUrl)

		res, err := client.Do(req)
		if err != nil {
			logger.Error("Request failed", "error", err.Error())
			return false
		}

		if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices {
			logger.Info("Request successful", "status", res.StatusCode)
			return true
		}

		logger.Error("Request unsuccessful", "status", res.StatusCode)
		return false
	}

	return pathHandler, nil
}
