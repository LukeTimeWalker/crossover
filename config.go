package main

import (
	"errors"
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

func DefaultConfig() *Config {
	return &Config{
		Port:     9999,
		Timeout:  15,
		LogLevel: 1,
	}
}

func GetConfig(filepath string) (*Config, error) {
	config := DefaultConfig()

	file, err := os.ReadFile(filepath)

	if err != nil {
		return nil, errors.New("error reading config file")
	}

	err = yaml.Unmarshal(file, config)

	if err != nil {
		return nil, errors.New("error unmarshaling config file")
	}

	return config, nil
}
