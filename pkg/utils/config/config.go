package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var Global Config

type Config struct {
	Image ImageStorage `yaml:"image"`
	JWT   JWT          `yaml:"jwt"`
	Mysql Mysql        `yaml:"mysql"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
	Database string `yaml:"database"`
}

type JWT struct {
	Secret     string        `yaml:"secret"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxRefresh time.Duration `yaml:"maxRefresh"`
}

type ImageStorage struct {
	Path string `yaml:"path"`
}

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	cfg := Config{}

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	Global = cfg

	return nil
}
