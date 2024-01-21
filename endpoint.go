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

func RegisterEndpoint(router *mux.Router, endpoint *Endpoint, client *http.Client, envs *map[string]interface{}) {
	logger := GetLogger()

	variables := make(Variables)

	variables["envs"] = *envs

	functions := make([]func(writer http.ResponseWriter, request *http.Request), len(endpoint.Sender))

	for index, sender := range endpoint.Sender {
		functions[index] = CreateSender(&sender, endpoint, variables, client)
	}

	handler := func(writer http.ResponseWriter, request *http.Request) {
		for index, function := range functions {
			logger.Debug("executing sender", "endpoint", endpoint.Receiver.Path, "sender", index)
			function(writer, request)
		}
	}

	path := path.Clean(fmt.Sprintf("/%s", endpoint.Receiver.Path))

	router.HandleFunc(path, handler).Methods(endpoint.Receiver.Method)
}

func CreateSender(sender *Sender, endpoint *Endpoint, variables Variables, client *http.Client) func(writer http.ResponseWriter, request *http.Request) {
	logger := GetLogger()

	var bodyTemplate *template.Template

	if sender.Body != "" {
		bodyTemplate = createTemplate("bodyTemplate", sender.Body)
	}

	var headerTemplates []*template.Template

	for key, value := range sender.Header {
		headerTemplates = append(headerTemplates, createTemplate(key, value))
	}

	var conditionTemplates []*template.Template

	for index, condition := range sender.Conditions {
		conditionTemplates = append(conditionTemplates, createTemplate(fmt.Sprintf("condition[%d]", index), condition))
	}

	urlTemplate := createTemplate("urlTemplate", sender.Url)

	handler := func(w http.ResponseWriter, request *http.Request) {
		if endpoint.Receiver.Body != nil {
			variables["body"] = LoadBody(request)
		}

		if endpoint.Receiver.UrlParams != nil {
			variables["url"] = LoadUrlParams(request.RequestURI, endpoint.Receiver.UrlParams)
		}

		senderBody := executeTemplate(bodyTemplate, variables)

		senderUrl := executeTemplate(urlTemplate, variables)

		// Send

		if !EvaluateConditions(conditionTemplates, variables) {
			logger.Info("Request was not send due the result of condtions evaluation")
			w.WriteHeader(http.StatusOK)
			return
		}

		req, err := http.NewRequest(sender.Method, senderUrl, bytes.NewReader([]byte(senderBody)))

		if err != nil {
			logger.Error("Could not create request", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if sender.Header != nil {
			for _, tmpl := range headerTemplates {
				req.Header.Add(tmpl.Name(), executeTemplate(tmpl, variables))
			}
		}

		logger.Info("Sending request", "method", sender.Method, "url", senderUrl)

		res, err := client.Do(req)

		if err != nil {
			logger.Error("request failed", "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		if res.StatusCode < http.StatusMultipleChoices && res.StatusCode >= http.StatusOK {
			logger.Info("Request successfull", "status", res.StatusCode)
		} else {
			logger.Error("request unsuccessfull", "status", res.StatusCode)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	return handler

}

func executeTemplate(template *template.Template, values interface{}) string {
	logger := GetLogger()

	var output bytes.Buffer

	err := template.Execute(&output, values)

	if err != nil {
		logger.Error("Could not execute template", "name", template.Name(), "error", err.Error())
	}

	return output.String()
}

func createTemplate(name string, templateText string) *template.Template {
	logger := GetLogger()

	template, err := template.New(name).Parse(templateText)

	if err != nil {
		logger.Error("Could not parse template", "name", name, "error", err.Error())
	}

	return template
}
