package utils_test

import (
	"os"
	"testing"

	"github.com/LukeTimeWalker/crossover/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnvsWithOneEnvSet(t *testing.T) {
	os.Setenv("TEST_ENV_2", "")
	os.Setenv("TEST_ENV_1", "1")

	environmenVariables := []string{"TEST_ENV_1"}
	envs := utils.LoadEnvs(environmenVariables)

	assert.Len(t, envs, 1)
	assert.Equal(t, envs["TEST_ENV_1"], "1")
}

func TestLoadEnvsWithOneEnvSetAndOneEnvEmpty(t *testing.T) {
	os.Setenv("TEST_ENV_2", "")
	os.Setenv("TEST_ENV_1", "1")

	environmenVariables := []string{"TEST_ENV_1", "TEST_ENV_2"}
	envs := utils.LoadEnvs(environmenVariables)

	assert.Len(t, envs, 2)
	assert.Equal(t, envs["TEST_ENV_1"], "1")
	assert.Equal(t, envs["TEST_ENV_1"], "1")
}

func TestLoadEnvThatIsNotSet(t *testing.T) {
	os.Setenv("TEST_ENV_2", "")
	os.Setenv("TEST_ENV_1", "1")

	environmenVariables := []string{"TEST_ENV_3"}
	envs := utils.LoadEnvs(environmenVariables)

	assert.Len(t, envs, 0)
	assert.Nil(t, envs["TEST_ENV_3"])
}
