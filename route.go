package main

import (
	"bytes"
	"errors"
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

	routeHandler := CreateRouteHandler(route.Receiver, variables, pathHandler)

	receiverEndpoint := path.Clean(fmt.Sprintf("/%s", route.Receiver.Endpoint))

	router.HandleFunc(receiverEndpoint, routeHandler).Methods(route.Receiver.Method)
}

func CreateRouteHandler(receiver Receiver, variables Variables, pathHandler []PathHandler) func(writer http.ResponseWriter, request *http.Request) {
	handler := func(writer http.ResponseWriter, request *http.Request) {

		if receiver.Body != nil {
			variables["body"] = LoadBody(request)
		}

		if receiver.UrlParams != nil {
			variables["url"] = LoadUrlParams(request.RequestURI, receiver.UrlParams)
		}

		statusCode := http.StatusOK
		for _, function := range pathHandler {
			success := function(variables)
			if !success {
				statusCode = 400
			}
		}
		writer.WriteHeader(statusCode)
	}
	return handler
}

func CreatePathHandler(receiver *Receiver, sender *Sender, client *http.Client) PathHandler {
	logger := GetLogger()

	var bodyTemplate *template.Template
	var err error

	if sender.Body != "" {
		bodyTemplate, err = createTemplate("bodyTemplate", sender.Body)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	var headerTemplates []*template.Template

	for key, value := range sender.Header {
		template, err := createTemplate(key, value)
		if err != nil {
			logger.Error(err.Error())
		}
		headerTemplates = append(headerTemplates, template)
	}

	var conditionTemplates []*template.Template

	for index, condition := range sender.Conditions {
		template, err := createTemplate(fmt.Sprintf("condition[%d]", index), condition)
		if err != nil {
			logger.Error(err.Error())
		}
		conditionTemplates = append(conditionTemplates, template)
	}

	urlTemplate, err := createTemplate("urlTemplate", sender.Endpoint)

	if err != nil {
		logger.Error(err.Error())
	}

	pathHandler := func(variables Variables) bool {
		senderBody, err := executeTemplate(bodyTemplate, variables)

		if err != nil {
			logger.Error(err.Error())
		}

		senderUrl, err := executeTemplate(urlTemplate, variables)

		if err != nil {
			logger.Error(err.Error())
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
				header, err := executeTemplate(tmpl, variables)

				if err != nil {
					logger.Error(err.Error())
				}

				req.Header.Add(tmpl.Name(), header)
			}
		}

		logger.Info("Sending request", "method", sender.Method, "url", senderUrl)

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

	return pathHandler

}

func executeTemplate(template *template.Template, values interface{}) (string, error) {
	var output bytes.Buffer

	err := template.Execute(&output, values)

	if err != nil {
		return "", errors.New("could not execute template")
		// logger.Error("Could not execute template", "name", template.Name(), "error", err.Error())
	}

	return output.String(), nil
}

func createTemplate(name string, templateText string) (*template.Template, error) {
	template, err := template.New(name).Parse(templateText)

	if err != nil {
		return nil, errors.New("could not parse template")
		// logger.Error("Could not parse template", "name", name, "error", err.Error())
	}

	return template, nil
}
