package conf

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fntsky/ddl_guard/internal/base/path"
	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "configs/config.yaml"

type Config struct {
	Server    ServerConfig   `yaml:"server"`
	Data      DataConfig     `yaml:"data"`
	VISUAL_AI VisualAIConfig `yaml:"VISUAL_AI"`
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

type VisualAIConfig struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
	Provider string `yaml:"provider"`
}

var (
	globalConfig     *Config
	globalConfigErr  error
	globalConfigOnce sync.Once
)

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

func LoadGlobal(configPath string) (*Config, error) {
	globalConfigOnce.Do(func() {
		globalConfig, globalConfigErr = ReadConfig(configPath)
	})
	return globalConfig, globalConfigErr
}

func Global() *Config {
	return globalConfig
}
