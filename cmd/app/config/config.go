package config

import (
	"fmt"
	"strings"

	"oats-docker/pkg/types"
)

type Config struct {
	Default DefaultOptions `yaml:"default"`
	Pgsql   PgsqlOptions   `yaml:"pgsql"`
	Cicd    CicdOptions    `yaml:"cicd"`
}

type DefaultOptions struct {
	Listen   int    `yaml:"listen"`
	LogType  string `yaml:"log_type"`
	LogDir   string `yaml:"log_dir"`
	LogLevel string `yaml:"log_level"`
	JWTKey   string `yaml:"jwt_key"`
}

type PgsqlOptions struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

type CicdOptions struct {
	Driver  string          `yaml:"driver"`
	Jenkins *JenkinsOptions `yaml:"jenkins"`
}

type JenkinsOptions struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (c *Config) Valid() error {
	if strings.ToLower(c.Default.LogType) == "file" {
		if len(c.Default.LogDir) == 0 {
			return fmt.Errorf("log_dir should be config when log type is file")
		}
	}

	switch c.Cicd.Driver {
	case "", types.Jenkins:
		j := c.Cicd.Jenkins
		if j == nil {
			return fmt.Errorf("jenkins config option missing")
		}
	default:
		return fmt.Errorf("unsupported cicd type %s", c.Cicd.Driver)
	}

	return nil
}
