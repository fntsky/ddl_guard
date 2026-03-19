package conf

import (
	"os"
	"path/filepath"

	"github.com/fntsky/ddl_guard/internal/base/path"
	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "configs/config.yaml"

type Config struct {
	Server ServerConfig `yaml:"server"`
	Data   DataConfig   `yaml:"data"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

type DataConfig struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Driver     string `yaml:"driver"`
	Connection string `yaml:"connection"`
}

func ReadConfig(configPath string) (*Config, error) {
	if len(configPath) == 0 {
		configPath = filepath.Join(path.ConfigFileDir, path.DefaultConfigFileName)
	}
	c := &Config{}
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(content, c); err != nil {
		return nil, err
	}
	return c, nil
}
