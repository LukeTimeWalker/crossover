package main

import (
	"os"
)

// Loads a defined list of environment variables
func LoadEnvs(environmentVariables []string) map[string]interface{} {
	envs := make(map[string]interface{})
	for _, env := range environmentVariables {
		loadEnv(env, envs)
	}
	return envs
}

// Loads a specific environment variable
func loadEnv(name string, vars map[string]interface{}) {
	logger := GetLogger()

	env, exists := os.LookupEnv(name)
	if exists {
		if env == "" {
			logger.Warn("Environment variable is set but empty", "variable", name)
			vars[name] = env
		} else {
			logger.Debug("Environment variable found", "variable", name)
			vars[name] = env
		}
	} else {
		logger.Warn("Environment variable not found", "variable", name)
	}
}
