package core

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
	Name       string   `yaml:"name"`
	Receiver   Receiver `yaml:"receiver"`
	Sender     []Sender `yaml:"sender"`
	Conditions []string `yaml:"conditions"`
}

type Receiver struct {
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"method"`
}

type Sender struct {
	Endpoint   string            `yaml:"endpoint"`
	Method     string            `yaml:"method"`
	Body       string            `yaml:"body"`
	Header     map[string]string `yaml:"header"`
	Client     Client            `yaml:"client"`
	Conditions []string          `yaml:"conditions"`
}

type Client struct {
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
