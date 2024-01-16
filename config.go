package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                  int        `yaml:"port"`
	Timeout               int        `yaml:"timeout"`
	LogLevel              int        `yaml:"logLevel"`
	Endpoints             []Endpoint `yaml:"endpoints"`
	EnvironementVariables []string   `yaml:"envs"`
}

type Endpoint struct {
	Receiver Receiver `yaml:"receiver"`
	Sender   []Sender `yaml:"sender"`
}

type Receiver struct {
	Path      string                 `yaml:"path"`
	Method    string                 `yaml:"method"`
	Body      map[string]interface{} `yaml:"body"`
	UrlParams map[string]string      `yaml:"urlParams"`
}

type Sender struct {
	Url        string            `yaml:"url"`
	Method     string            `yaml:"method"`
	Body       string            `yaml:"body"`
	Header     map[string]string `yaml:"header"`
	Conditions []string          `yaml:"conditions"`
}

func GetConfig(filepath string) *Config {
	config := Config{}

	logger := GetLogger()

	file, err := os.ReadFile(filepath)

	if err != nil {
		logger.Error("error reading config file", "file", filepath, "error", err)
	}

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		logger.Error("error unmarshaling config file", "file", filepath, "error", err)
	}

	return &config
}
