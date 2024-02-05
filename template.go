package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func CreateTemplate(name string, templateText string) (*template.Template, error) {
	template, err := template.New(name).Parse(templateText)

	if err != nil {
		return nil, fmt.Errorf("could not parse template name=%s error=%s", name, err.Error())
	}

	return template, nil
}

func ExecuteTemplate(template *template.Template, values interface{}) (string, error) {
	var output bytes.Buffer

	err := template.Execute(&output, values)

	if err != nil {
		return "", fmt.Errorf("could not execute template name=%s error=%s", template.Name(), err.Error())
	}

	return output.String(), nil
}
