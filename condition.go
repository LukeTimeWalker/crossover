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

func EvaluateConditions(conditions []*template.Template, variables Variables) bool {
	for _, condition := range conditions {
		if !EvaluateCondition(condition, variables) {
			return false
		}
	}
	return true
}

func EvaluateCondition(condition *template.Template, variables Variables) bool {
	logger := GetLogger()

	value, err := ExecuteTemplate(condition, variables)

	if err != nil {
		logger.Error(err.Error())
	}

	result, err := strconv.ParseBool(value)

	if err != nil {
		logger.Error("condition does not evaluate to a boolean value", "condition", condition.Name(), "value", value)
		return false
	}
	logger.Debug("condition evaluates to boolean value", "condition", condition.Name(), "value", result)

	return result
}
