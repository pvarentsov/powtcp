package config

import (
	"flag"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// Config - config structure
type Config struct {
	Server struct {
		LogLevel          int    `yaml:"log_level"`
		LogJson           bool   `yaml:"log_json"`
		Address           string `yaml:"address"`
		ShutdownTimeout   int    `yaml:"shutdown_timeout"`
		ConnectionTimeout int    `yaml:"connection_timeout"`
	} `yaml:"server"`

	Client struct {
		LogLevel      int    `yaml:"log_level"`
		LogJson       bool   `yaml:"log_json"`
		ServerAddress string `yaml:"server_address"`
	} `yaml:"client"`

	Hashcash struct {
		Bits               int `yaml:"bits"`
		ComputeMaxAttempts int `yaml:"compute_max_attempts"`
		TTL                int `yaml:"ttl"`
	} `yaml:"hashcash"`
}

// ParseByFlag - parse config from file by flag
func ParseByFlag(flagName string) (config *Config, err error) {
	var path string

	flag.StringVar(&path, flagName, "", "")
	flag.Parse()

	if path == "" {
		return nil, fmt.Errorf("config path not set, pass config path using '--%s' flag", flagName)
	}
	if config, err = Parse(path); err != nil {
		return nil, err
	}

	return
}

// Parse - parse config from file
func Parse(dest string) (config *Config, err error) {
	f, err := os.Open(path.Clean(dest))
	if err != nil {
		return nil, err
	}

	defer f.Close()

	if err = yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return
}
