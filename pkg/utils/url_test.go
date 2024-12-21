package utils_test

import (
	"testing"

	"github.com/LukeTimeWalker/crossover/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoadUrlParams(t *testing.T) {
	u := "https://example.com/?key1=value1&key2=value2&key3=value3"
	expectedParams := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
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
