package config

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ListenPort string   `json:"listen_port" yaml:"listen_port"`
	Backends   []string `json:"backends" yaml:"backends"`
}

// Загружает конфигурацию из JSON файла
func LoadConfigFromJSON(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}

// Загружает конфигурацию из YAML файла
func LoadConfigFromYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	return &config, err
}
