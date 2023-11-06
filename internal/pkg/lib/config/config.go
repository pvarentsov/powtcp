package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - config structure
type Config struct {
	Server   `yaml:"server" env-prefix:"SERVER_"`
	Client   `yaml:"client" env-prefix:"CLIENT_"`
	Hashcash `yaml:"hashcash" env-prefix:"HASHCASH_"`
}

// Server - server config structure
type Server struct {
	LogLevel          int    `yaml:"log_level" env:"LOG_LEVEL" env-default:"0"`
	LogJson           bool   `yaml:"log_json" env:"LOG_JSON" env-default:"false"`
	Address           string `yaml:"address" env:"ADDRESS" env-default:":8080"`
	ShutdownTimeout   int    `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"1000"`
	ConnectionTimeout int    `yaml:"connection_timeout" env:"CONNECTION_TIMEOUT" env-default:"30000"`
}

// Client - client config structure
type Client struct {
	LogLevel      int    `yaml:"log_level" env:"LOG_LEVEL" env-default:"0"`
	LogJson       bool   `yaml:"log_json" env:"LOG_JSON" env-default:"false"`
	ServerAddress string `yaml:"server_address" env:"SERVER_ADDRESS" env-default:":8080"`
}

// Hashcash - Hashcash config structure
type Hashcash struct {
	Bits               int `yaml:"bits" env:"BITS" env-default:"5"`
	ComputeMaxAttempts int `yaml:"compute_max_attempts"  env:"COMPUTE_MAX_ATTEMPTS" env-default:"100000000"`
	TTL                int `yaml:"ttl"  env:"TTL" env-default:"60000"`
}

// Parse - parse config from file by flag or from env or use default
func Parse(flagName string) (config *Config, err error) {
	var path string

	flag.StringVar(&path, flagName, "", "")
	flag.Parse()

	if path == "" {
		return ParseFromEnv()
	}

	return ParseFromFile(path)
}

// ParseFromEnv - parse config from environment variables
func ParseFromEnv() (*Config, error) {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// ParseFromFile - parse config from file
func ParseFromFile(path string) (*Config, error) {
	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
