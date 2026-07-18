package utils_test

import (
	"testing"

	"github.com/LukeTimeWalker/crossover/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoadUrlParams(t *testing.T) {
	u := "https://example.com/?key1=value1&key2=value2&key3=value3"
	expectedParams := map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	params, err := utils.LoadUrlParams(u)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, params)
}

func TestLoadUrlParamsWithMultipleValuesPerKey(t *testing.T) {
	u := "https://example.com/?key1=value1&key2=value2&key3=value3&key3=value4&key1=value5"
	expectedParams := map[string]any{
		"key1": []string{"value1", "value5"},
		"key2": "value2",
		"key3": []string{"value3", "value4"},
	}

	params, err := utils.LoadUrlParams(u)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, params)
}

func TestLoadUrlParamsWithAMissingValue(t *testing.T) {
	u := "https://example.com/?key1&key2=value2"
	expectedParams := map[string]any{
		"key1": "",
		"key2": "value2",
	}

	params, err := utils.LoadUrlParams(u)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, params)
}

func TestLoadUrlParamsWithParsingErrors(t *testing.T) {
	u := "https://example.com/?key1=value1&key2=value2&key3=value3%ZZ"

	params, err := utils.LoadUrlParams(u)

	assert.Error(t, err)
	assert.Nil(t, params)
}
