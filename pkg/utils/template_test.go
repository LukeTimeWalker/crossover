package utils_test

import (
	"testing"

	"github.com/LukeTimeWalker/crossover/pkg/core"
	"github.com/LukeTimeWalker/crossover/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const testTemplate = `
		{
			"title": "{{ .env.TITLE }}",
			"description": "Body: {{ .body.test }}  \n{{ len .body.array }}  \n{{ .url.url }}"
		}
	`

func TestCreatingAValidTemplate(t *testing.T) {

	template, err := utils.CreateTemplate("testTemplate", "'Bearer {{ .env.TEST_TOKEN }}'")
	assert.NoError(t, err)
	assert.NotNil(t, template)
}

func TestCreatingAnInvalidTemplate(t *testing.T) {
	template, err := utils.CreateTemplate("testTemplate", "'Bearer {{ .env.TEST_TOKEN }'")
	assert.Error(t, err)
	assert.Nil(t, template)
}

func TestExecutingATemplateWithAllValues(t *testing.T) {
	template, _ := utils.CreateTemplate("testTemplate", testTemplate)
	values := core.Variables{
		"env":  map[string]any{"TITLE": "testTitle"},
		"body": map[string]any{"test": "bodyContent", "array": []string{"1", "2", "3"}},
		"url":  map[string]any{"url": "param"},
	}
	result, err := utils.ExecuteTemplate(template, values)
	assert.NoError(t, err)
	assert.Equal(t, `
		{
			"title": "testTitle",
			"description": "Body: bodyContent  \n3  \nparam"
		}
	`, result)
}

func TestExecutingTemplateWithMissingValues(t *testing.T) {
	template, _ := utils.CreateTemplate("testTemplate", testTemplate)
	values := core.Variables{
		"env":  map[string]any{"TITLE": "testTitle"},
		"body": map[string]any{"test": "bodyContent"},
		"url":  map[string]any{"url": "param"}}
	result, err := utils.ExecuteTemplate(template, values)
	assert.Error(t, err)
	assert.Empty(t, result)
}
