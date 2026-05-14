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
	WECHAT    WechatConfig   `yaml:"WECHAT"`
	JWT       JWTConfig      `yaml:"jwt"`
	Publish   PublishConfig  `yaml:"publish"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr" env:"SERVER_HTTP_ADDR"`
}

type DataConfig struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Driver     string `yaml:"driver" env:"DATABASE_DRIVER"`
	Connection string `yaml:"connection" env:"DATABASE_CONNECTION"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env:"REDIS_ADDR"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

type VisualAIConfig struct {
	APIKey   string `yaml:"api_key" env:"VISUAL_AI_API_KEY"`
	Endpoint string `yaml:"endpoint" env:"VISUAL_AI_ENDPOINT"`
	Model    string `yaml:"model" env:"VISUAL_AI_MODEL"`
	Provider string `yaml:"provider" env:"VISUAL_AI_PROVIDER"`
}

type EmailOTPConfig struct {
	SMTP SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host     string `yaml:"host" env:"SMTP_HOST"`
	Port     int    `yaml:"port" env:"SMTP_PORT"`
	Username string `yaml:"username" env:"SMTP_USERNAME"`
	Password string `yaml:"password" env:"SMTP_PASSWORD"`
}

type JWTConfig struct {
	Secret           string `yaml:"secret" env:"JWT_SECRET"`
	AccessTTLMinutes int    `yaml:"access_ttl_minutes" env:"JWT_ACCESS_TTL_MINUTES"`
	RefreshTTLHours  int    `yaml:"refresh_ttl_hours" env:"JWT_REFRESH_TTL_HOURS"`
}

type WechatConfig struct {
	AppID     string `yaml:"app_id" env:"WECHAT_APP_ID"`
	AppSecret string `yaml:"app_secret" env:"WECHAT_APP_SECRET"`
}

// PublishConfig 推送配置
type PublishConfig struct {
	Email EmailPublishConfig `yaml:"email"`
}

// EmailPublishConfig 邮件推送配置
type EmailPublishConfig struct {
	Enabled bool       `yaml:"enabled" env:"PUBLISH_EMAIL_ENABLED"`
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
	// 用环境变量覆盖配置
	applyEnvVars(c)
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
