package config

import (
	"os"

	"github.com/0chain/bandwidth_marketplace/code/core/errors"
	"gopkg.in/yaml.v2"
)

type (
	// Config represents configs stores in Path.
	Config struct {
		Port              int    `yaml:"port"`
		GRPCAddress       string `yaml:"grpc_address"`
		GRPCServerTimeout int    `yaml:"grpc_server_timeout"` // in seconds

		ConsumerAddress string `yaml:"consumer_address"`
		ProviderAddress string `yaml:"provider_address"`

		Handler Handler `yaml:"handler"`
		HSS     HSS     `yaml:"hss"`
		Logging Logging `yaml:"logging"`
	}

	// Handler represents config options described in "handler" section of the config yaml file.
	// Handler must be a field of Config struct
	Handler struct {
		RateLimit float64 `yaml:"rate_limit"` // per second
	}

	// Logging represents config options described in "logging" section of the config yaml file.
	// Database must be a field of Config struct
	Logging struct {
		Level string `yaml:"level"`
		Dir   string `yaml:"dir"`
	}

	// HSS represents config options described in "hss" section of the config yaml file.
	// HSS must be a field of Config struct
	HSS struct {
		Users []string `yaml:"users"`
	}
)

const (
	// Path is a constant stores path to config file from root application directory.
	Path = "./config/magma-config.yaml"
)

// Read reads configs from config file existing in Path.
//
// If an error occurs during execution, the program terminates with code 2 and the error will be written in os.Stderr.
//
// Read should be used only once while application is starting.
func Read() *Config {
	f, err := os.Open(Path)
	if err != nil {
		errors.ExitErr("err while open config file", err, 2)
	}
	defer func(f *os.File) { _ = f.Close() }(f)

	decoder := yaml.NewDecoder(f)
	cfg := new(Config)
	err = decoder.Decode(cfg)
	if err != nil {
		errors.ExitErr("err while decoding config file", err, 2)
	}

	return cfg
}
