package main

import (
	"os"
)

func LoadEnvs(environmentVariables []string) map[string]interface{} {
	envs := make(map[string]interface{})
	for _, env := range environmentVariables {
		loadEnv(env, envs)
	}
	return envs
}

func loadEnv(name string, vars map[string]interface{}) {
	logger := GetLogger()

	env, exists := os.LookupEnv(name)

	if exists {
		if env == "" {
			logger.Warn("Environment variable is empty", "variable", name)
			vars[name] = env
		} else {
			logger.Debug("Environment variable found", "variable", name)
			vars[name] = env
		}
	} else {
		logger.Warn("Environment variable not found", "variable", name)
	}
}
