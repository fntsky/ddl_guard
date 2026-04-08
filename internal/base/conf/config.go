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
	Redis     RedisConfig    `yaml:"redis"`
	VISUAL_AI VisualAIConfig `yaml:"VISUAL_AI"`
	EMAIL_OTP EmailOTPConfig `yaml:"EMAIL_OTP"`
	JWT       JWTConfig      `yaml:"jwt"`
	Publish   PublishConfig  `yaml:"publish"`
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

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type VisualAIConfig struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
	Provider string `yaml:"provider"`
}

type EmailOTPConfig struct {
	SMTP SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type JWTConfig struct {
	Secret           string `yaml:"secret"`
	AccessTTLMinutes int    `yaml:"access_ttl_minutes"`
	RefreshTTLHours  int    `yaml:"refresh_ttl_hours"`
}

// PublishConfig 推送配置
type PublishConfig struct {
	Email EmailPublishConfig `yaml:"email"`
}

// EmailPublishConfig 邮件推送配置
type EmailPublishConfig struct {
	Enabled bool       `yaml:"enabled"`
	SMTP    SMTPConfig `yaml:"smtp"`
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
