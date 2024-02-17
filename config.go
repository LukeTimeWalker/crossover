package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                  int           `yaml:"port"`
	Timeout               int           `yaml:"timeout"`
	LogLevel              int           `yaml:"logLevel"`
	Routes                []RouteConfig `yaml:"routes"`
	EnvironementVariables []string      `yaml:"envs"`
}

type RouteConfig struct {
	Receiver   ReceiverConfig `yaml:"receiver"`
	Sender     []SenderConfig `yaml:"sender"`
	Conditions []string       `yaml:"conditions"`
}

type ReceiverConfig struct {
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"method"`
}

type SenderConfig struct {
	Endpoint   string            `yaml:"endpoint"`
	Method     string            `yaml:"method"`
	Body       string            `yaml:"body"`
	Header     map[string]string `yaml:"header"`
	Client     ClientConfig      `yaml:"client"`
	Conditions []string          `yaml:"conditions"`
}

type ClientConfig struct {
	Timeout       int  `yaml:"timeout"`
	TlsSkipVerify bool `yaml:"tlsSkipVerify"`
}

func DefaultConfig() *Config {
	return &Config{
		Port:     9999,
		Timeout:  10,
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
