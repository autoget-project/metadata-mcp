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
	ThePornDBAPIToken    string `yaml:"theporndb_api_token"`
	MetaTubeAPIURL       string `yaml:"metatube_api_url"`
	MetaTubeAPIKEY       string `yaml:"metatube_api_key"`
	WikipediaLanguage    string `yaml:"wikipedia_language"`
	JAVActorJSONFilepath string `yaml:"jav_actor_json_filepath"`
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

	if c.ThePornDBAPIToken == "" {
		return fmt.Errorf("ThePornDB_API_KEY is required")
	}
	if c.MetaTubeAPIURL == "" {
		return fmt.Errorf("MetaTube_API_URL is required")
	}
	// MetaTube_API_KEY is optional

	if c.WikipediaLanguage == "" {
		// default language is zh
		c.WikipediaLanguage = "zh"
	}
	if c.JAVActorJSONFilepath == "" {
		return fmt.Errorf("JAV_ACTOR_JSON_FILEPATH is required")
	}
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

func ReadConfigFromEnv() (*Config, error) {
	conf := &Config{}

	if portStr := os.Getenv("PORT"); portStr != "" {
		_, err := fmt.Sscanf(portStr, "%d", &conf.Port)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT environment variable: %w", err)
		}
	}
	conf.TMDBAPIKey = os.Getenv("TMDB_API_KEY")
	conf.TMDBResponseLanguage = os.Getenv("TMDB_RESPONSE_LANGUAGE")
	conf.ThePornDBAPIToken = os.Getenv("TPDB_API_TOKEN")
	conf.MetaTubeAPIURL = os.Getenv("METATUBE_API_URL")
	conf.MetaTubeAPIKEY = os.Getenv("METATUBE_API_KEY")
	conf.WikipediaLanguage = os.Getenv("WIKIPEDIA_LANGUAGE")
	conf.JAVActorJSONFilepath = os.Getenv("JAV_ACTOR_JSON_FILEPATH")

	err := conf.validate()
	if err != nil {
		return nil, err
	}
	return conf, nil
}
