package config

import (
	"os"
	"time"

	"github.com/DaryaPe/hw-test/hw12_13_14_15_calendar/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config конфигурация сервиса.
type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// LoggerConfig конфигурация для логгера.
type LoggerConfig struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

// ServerConfig конфигурация для HTTP сервера.
type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

// DatabaseConfig конфигурация для базы данных.
type DatabaseConfig struct {
	Source   string        `yaml:"source"`
	Username string        `yaml:"user"`
	Password string        `yaml:"pass"`
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Database string        `yaml:"database"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Apply применяет значение из конфигурационного файла.
func (c *Config) Apply(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f) // nolint
	if err = decoder.Decode(c); err != nil {
		return errors.Wrap(err, "decoder.Decode")
	}
	return nil
}

func New() *Config {
	return &Config{}
}
