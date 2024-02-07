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

	variables := make(Variables)

	variables["envs"] = *envs

	pathHandler := make([]PathHandler, len(route.Sender))

	for index, sender := range route.Sender {
		handler, err := CreatePathHandler(&route.Receiver, &sender)
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

		body, err := LoadBody(request)

		if err != nil {
			logger.Error(err.Error())
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		variables["body"] = body

		if route.Receiver.UrlParams != nil {
			urlParams, err := LoadUrlParams(request.RequestURI, route.Receiver.UrlParams)

			if err != nil {
				logger.Error(err.Error())
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			variables["url"] = urlParams
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

func CreatePathHandler(receiverConfig *ReceiverConfig, senderConfig *SenderConfig) (PathHandler, error) {

	var bodyTemplate *template.Template
	var err error

	if senderConfig.Body != "" {
		bodyTemplate, err = CreateTemplate("bodyTemplate", senderConfig.Body)
		if err != nil {
			return nil, err
		}
	}

	var headerTemplates []*template.Template

	for key, value := range senderConfig.Header {
		template, err := CreateTemplate(key, value)
		if err != nil {
			return nil, err
		}
		headerTemplates = append(headerTemplates, template)
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
			logger.Info("Request was not send due the result of condtions evaluation")
			return false
		}

		req, err := http.NewRequest(senderConfig.Method, senderUrl, bytes.NewReader([]byte(senderBody)))

		if err != nil {
			logger.Error("Could not create request", "error", err.Error())
			return false
		}

		if senderConfig.Header != nil {
			for _, tmpl := range headerTemplates {
				header, err := ExecuteTemplate(tmpl, variables)

				if err != nil {
					logger.Error(err.Error())
				}

				req.Header.Add(tmpl.Name(), header)
			}
		}

		logger.Debug("Sending request", "method", senderConfig.Method, "url", senderUrl)

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
