package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                  int      `yaml:"port"`
	Timeout               int      `yaml:"timeout"`
	LogLevel              int      `yaml:"logLevel"`
	Routes                []Route  `yaml:"routes"`
	EnvironementVariables []string `yaml:"envs"`
}

type Route struct {
	Receiver   Receiver `yaml:"receiver"`
	Sender     []Sender `yaml:"sender"`
	Conditions []string `yaml:"contitions"`
}

type Receiver struct {
	Endpoint  string                 `yaml:"endpoint"`
	Method    string                 `yaml:"method"`
	Body      map[string]interface{} `yaml:"body"`
	UrlParams map[string]string      `yaml:"urlParams"`
}

type Sender struct {
	Endpoint   string            `yaml:"endpoint"`
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
		return nil, fmt.Errorf("error reading config file error=%s", err.Error())
	}

	err = yaml.Unmarshal(file, config)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config file error=%s", err.Error())
	}

	return config, nil
}
