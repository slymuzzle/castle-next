package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type ServerConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type AuthConfig struct {
	Secret string `koanf:"secret"`
}

type NatsConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type DatabaseConfig struct {
	Database string `koanf:"database"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Driver   string `koanf:"driver"`
}

type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Auth     AuthConfig     `koanf:"auth"`
	Nats     NatsConfig     `koanf:"nats"`
	Database DatabaseConfig `koanf:"database"`
}

var (
	configFilePath        = "configs/server.yml"
	envConfigPrefix       = "JOURNEYHUB_"
	ErrLoadingConfig      = "error loading config: %v"
	ErrUnmarshalingConfig = "error unmarshaling config: %v"
)

type Service interface {
	LoadConfig() (Config, error)
}

type service struct {
	koanf koanf.Koanf
}

func NewService() Service {
	return &service{
		koanf: *koanf.New("."),
	}
}

func (s *service) LoadConfig() (Config, error) {
	if err := s.koanf.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf(ErrLoadingConfig, err)
	}

	s.koanf.Load(env.Provider(envConfigPrefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, envConfigPrefix)), "_", ".")
	}), nil)

	var config Config
	if err := s.koanf.Unmarshal("", &config); err != nil {
		return Config{}, fmt.Errorf(ErrUnmarshalingConfig, err)
	}

	return config, nil
}
