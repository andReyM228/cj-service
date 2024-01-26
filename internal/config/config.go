package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type (
	Config struct {
		TgBot TgBot `yaml:"tg-bot" validate:"required"`
		Extra Extra `yaml:"extra" validate:"required"`
	}

	TgBot struct {
		Token string `yaml:"token" validate:"required"`
	}

	GoogleSheets struct {
		UrlNames        string `yaml:"urlNames" validate:"required"`
		UrlQuestions    string `yaml:"urlQuestions" validate:"required"`
		UrlAllQuestions string `yaml:"urlAllQuestions" validate:"required"`
		UrlInfo         string `yaml:"urlInfo" validate:"required"`
	}

	Extra struct {
		GoogleSheets GoogleSheets `yaml:"google-sheets" validate:"required"`
	}
)

func ParseConfig() (Config, error) {
	file, err := os.ReadFile("./cmd/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg, nil
}
