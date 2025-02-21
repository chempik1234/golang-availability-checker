package parsed

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Webhook struct {
	Url     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers"`
}

type Checker struct {
	Name            string    `yaml:"name"`
	Host            string    `yaml:"host"`
	Protocol        string    `yaml:"protocol"`
	Port            string    `yaml:"port"`
	IntervalSeconds int       `yaml:"interval_seconds"`
	TimeoutSeconds  int       `yaml:"timeout_seconds"`
	Webhooks        []Webhook `yaml:"webhooks"`
}

type Service struct {
	Version  string    `yaml:"version"`
	Services []Checker `yaml:"services"`
}

func ParseService(filepath string) (*Service, error) {
	var s Service
	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
