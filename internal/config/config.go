package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/labstack/gommon/log"
)

type Config struct {
	MassiveAPIKey  string
	ScannerOptions struct {
		MinPrice               float64 `yaml:"minPrice"`
		MaxPrice               float64 `yaml:"maxPrice"`
		MinVolume              uint64  `yaml:"minVolume"`
		MinPremarketGapPercent float64 `yaml:"minPremarketGapPercent"`
		MaxPremarketGapPercent float64 `yaml:"maxPremarketGapPercent"`
	} `yaml:"scannerOptions"`
}

func LoadConfig() (*Config, error) {
	config := Config{}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(cwd, "config.yml"))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	config.MassiveAPIKey = os.Getenv("MASSIVE_API_KEY")
	log.Debugf("Config loaded: %+v", config)
	return &config, err
}
