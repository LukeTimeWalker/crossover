package main

import (
	"fmt"
	"strconv"
	"text/template"
)

func CreateConditions(conditions []string) ([]*template.Template, error) {
	conditionTemplates := make([]*template.Template, len(conditions))

	for index, condition := range conditions {
		template, err := CreateTemplate(fmt.Sprintf("condition[%d]", index), condition)

		if err != nil {
			return nil, err
		}

		conditionTemplates = append(conditionTemplates, template)
	}

	return conditionTemplates, nil
}

func EvaluateConditions(conditions []*template.Template, variables Variables) (bool, error) {
	for _, condition := range conditions {
		result, err := evaluateCondition(condition, variables)

		if err != nil {
			return false, err
		}

		if !result {
			return false, nil
		}
	}

	return true, nil
}

func evaluateCondition(condition *template.Template, variables Variables) (bool, error) {
	value, err := ExecuteTemplate(condition, variables)

	if err != nil {
		return false, err
	}

	result, err := strconv.ParseBool(value)

	if err != nil {
		return false, fmt.Errorf("condition does not evaluate to a boolean value condition=%s value=%s", condition.Name(), value)
	}

	return result, nil
}
