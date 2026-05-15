package core

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port                  int      `yaml:"port"`
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
	Timeout       string `yaml:"timeout"`
	TlsSkipVerify bool   `yaml:"tlsSkipVerify"`
}

func (s *Sender) withDefaults() {
	s.Client.withDefaults()
}

func (r *Route) withDefaults() {
	for i := range r.Sender {
		r.Sender[i].withDefaults()
	}
}

func (c *Config) withDefaults() {
	if c.Port == 0 {
		c.Port = 9999
	}
	if c.LogLevel == 0 {
		c.LogLevel = 1
	}
	for i := range c.Routes {
		c.Routes[i].withDefaults()
	}
}

func (c *Client) withDefaults() {
	if c.Timeout == "" {
		c.Timeout = "10s"
	}
}

func GetConfig(filepath string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file error=%s", err.Error())
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config file error=%s", err.Error())
	}

	config.withDefaults()

	return config, nil
}
