package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Port                 int    `yaml:"port"`
	TMDBAPIKey           string `yaml:"tmdb_api_key"`
	TMDBResponseLanguage string `yaml:"tmdb_response_language"`
	ThePornDBAPIKey      string `yaml:"theporndb_api_key"`
	MetaTubeAPIURL       string `yaml:"metatube_api_url"`
	MetaTubeAPIKEY       string `yaml:"metatube_api_key"`
}

func (c *Config) validate() error {
	if c.Port == 0 {
		// default port is 8080
		c.Port = 8080
	}

	if c.TMDBAPIKey == "" {
		return fmt.Errorf("TMDB_API_KEY is required")
	}
	if c.TMDBResponseLanguage == "" {
		// default language is zh-CN
		c.TMDBResponseLanguage = "zh-CN"
	}

	if c.ThePornDBAPIKey == "" {
		return fmt.Errorf("ThePornDB_API_KEY is required")
	}
	if c.MetaTubeAPIURL == "" {
		return fmt.Errorf("MetaTube_API_URL is required")
	}
	// MetaTube_API_KEY is optional
	return nil
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}
	err = conf.validate()
	if err != nil {
		return nil, err
	}
	return conf, nil
}
