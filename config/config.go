package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"emperror.dev/errors"
)

var (
	cfg Config
)

type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

type ServerConfig struct {
	Port     string        `json:"port"`
	Timeout  time.Duration `json:"timeout"`
	LogLevel string        `json:"logLevel"`
	BaseAddr string        `json:"baseAddr"`
}

type DatabaseConfig struct {
	ConnectionString string        `json:"conectionString"`
	MigrationsDir    string        `json:"migrationsDir"`
	MaxOpenConns     int           `json:"maxOpenConns"`
	MaxIdleConns     int           `json:"maxIdleConns"`
	DBTimeout        time.Duration `json:"dbTimeout"`
}

func Instance() Config {
	return cfg
}

func ReadConfig(env, dir string) (Config, error) {
	configPath := path.Join(dir, "config")

	err := cfg.loadJSON(path.Join(configPath, fmt.Sprintf("config.%s.json", env)))
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c *Config) loadJSON(filepath string) error {
	if _, err := os.Stat(filepath); err != nil {
		return errors.Wrap(err, "config file not found")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return errors.Wrapf(err, "error opening config file: %s", filepath)
	}

	byteValue, _ := io.ReadAll(file)

	if err = file.Close(); err != nil {
		return errors.Wrapf(err, "error closing config file %v", filepath)
	}

	err = json.Unmarshal(byteValue, c)
	if err != nil {
		return errors.Wrapf(err, "error unmarshalling config file %v", filepath)
	}

	return nil
}
