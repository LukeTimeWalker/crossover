package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func CreateTemplate(name string, templateText string) (*template.Template, error) {
	tmpl, err := template.New(name).Parse(templateText)

	if err != nil {
		return nil, fmt.Errorf("could not parse template name=%s error=%s", name, err.Error())
	}

	return tmpl, nil
}

func ExecuteTemplate(tmpl *template.Template, values interface{}) (string, error) {
	var output bytes.Buffer

	err := tmpl.Execute(&output, values)

	if err != nil {
		return "", fmt.Errorf("could not execute template name=%s error=%s", tmpl.Name(), err.Error())
	}

	return output.String(), nil
}
