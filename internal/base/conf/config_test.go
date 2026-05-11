package conf

import (
	"os"
	"reflect"
	"testing"
)

func TestApplyEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		initial  Config
		expected Config
	}{
		{
			name:    "override string field",
			envVars: map[string]string{"JWT_SECRET": "new_secret"},
			initial: Config{JWT: JWTConfig{Secret: "old"}},
			expected: Config{JWT: JWTConfig{
				Secret:           "new_secret",
				AccessTTLMinutes: 0,
				RefreshTTLHours:  0,
			}},
		},
		{
			name:    "override int field",
			envVars: map[string]string{"JWT_ACCESS_TTL_MINUTES": "30"},
			initial: Config{JWT: JWTConfig{AccessTTLMinutes: 15}},
			expected: Config{JWT: JWTConfig{
				Secret:           "",
				AccessTTLMinutes: 30,
				RefreshTTLHours:  0,
			}},
		},
		{
			name:    "override bool field",
			envVars: map[string]string{"PUBLISH_EMAIL_ENABLED": "false"},
			initial: Config{Publish: PublishConfig{Email: EmailPublishConfig{Enabled: true}}},
			expected: Config{Publish: PublishConfig{Email: EmailPublishConfig{
				Enabled: false,
			}}},
		},
		{
			name:    "override multiple fields",
			envVars: map[string]string{
				"REDIS_ADDR":     "redis:6379",
				"REDIS_PASSWORD": "newpass",
				"REDIS_DB":       "5",
			},
			initial: Config{Redis: RedisConfig{Addr: "localhost:6379", Password: "", DB: 0}},
			expected: Config{Redis: RedisConfig{
				Addr:     "redis:6379",
				Password: "newpass",
				DB:       5,
			}},
		},
		{
			name:    "nested struct - SMTP config",
			envVars: map[string]string{
				"SMTP_HOST":     "smtp.example.com",
				"SMTP_PORT":     "587",
				"SMTP_USERNAME": "user@example.com",
				"SMTP_PASSWORD": "newpass",
			},
			initial: Config{EMAIL_OTP: EmailOTPConfig{SMTP: SMTPConfig{
				Host:     "old.host.com",
				Port:     465,
				Username: "old@example.com",
				Password: "oldpass",
			}}},
			expected: Config{EMAIL_OTP: EmailOTPConfig{SMTP: SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user@example.com",
				Password: "newpass",
			}}},
		},
		{
			name:     "no env vars - keep original",
			envVars:  map[string]string{},
			initial:  Config{JWT: JWTConfig{Secret: "original", AccessTTLMinutes: 15}},
			expected: Config{JWT: JWTConfig{Secret: "original", AccessTTLMinutes: 15}},
		},
		{
			name:    "database config",
			envVars: map[string]string{
				"DATABASE_DRIVER":     "mysql",
				"DATABASE_CONNECTION": "user:pass@tcp(localhost:3306)/db",
			},
			initial: Config{Data: DataConfig{Database: DatabaseConfig{
				Driver:     "postgres",
				Connection: "postgres://localhost/db",
			}}},
			expected: Config{Data: DataConfig{Database: DatabaseConfig{
				Driver:     "mysql",
				Connection: "user:pass@tcp(localhost:3306)/db",
			}}},
		},
		{
			name:    "visual AI config",
			envVars: map[string]string{
				"VISUAL_AI_API_KEY":   "new-api-key",
				"VISUAL_AI_ENDPOINT":  "https://new.endpoint.com",
				"VISUAL_AI_MODEL":     "new-model",
				"VISUAL_AI_PROVIDER":  "openai",
			},
			initial: Config{VISUAL_AI: VisualAIConfig{
				APIKey:   "old-key",
				Endpoint: "https://old.endpoint.com",
				Model:    "old-model",
				Provider: "glm",
			}},
			expected: Config{VISUAL_AI: VisualAIConfig{
				APIKey:   "new-api-key",
				Endpoint: "https://new.endpoint.com",
				Model:    "new-model",
				Provider: "openai",
			}},
		},
		{
			name:    "server HTTP addr",
			envVars: map[string]string{"SERVER_HTTP_ADDR": "0.0.0.0:9090"},
			initial: Config{Server: ServerConfig{HTTP: HTTPConfig{Addr: "0.0.0.0:8080"}}},
			expected: Config{Server: ServerConfig{HTTP: HTTPConfig{
				Addr: "0.0.0.0:9090",
			}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			// 应用环境变量
			applyEnvVars(&tt.initial)

			// 验证结果
			if !configEqual(tt.initial, tt.expected) {
				t.Errorf("config mismatch\n got: %+v\n want: %+v", tt.initial, tt.expected)
			}
		})
	}
}

func TestSetFieldValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		kind     reflect.Kind
		expected interface{}
		hasError bool
	}{
		{"string value", "hello", reflect.String, "hello", false},
		{"int value", "42", reflect.Int, int64(42), false},
		{"bool true", "true", reflect.Bool, true, false},
		{"bool false", "false", reflect.Bool, false, false},
		{"invalid int", "abc", reflect.Int, nil, true},
		{"invalid bool", "maybe", reflect.Bool, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var field reflect.Value
			switch tt.kind {
			case reflect.String:
				var s string
				field = reflect.ValueOf(&s).Elem()
			case reflect.Int:
				var i int
				field = reflect.ValueOf(&i).Elem()
			case reflect.Bool:
				var b bool
				field = reflect.ValueOf(&b).Elem()
			}

			err := setFieldValue(field, tt.value)

			if tt.hasError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			var got interface{}
			switch tt.kind {
			case reflect.String:
				got = field.String()
			case reflect.Int:
				got = field.Int()
			case reflect.Bool:
				got = field.Bool()
			}

			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

// configEqual 比较两个 Config 是否相等
func configEqual(a, b Config) bool {
	return a.Server.HTTP.Addr == b.Server.HTTP.Addr &&
		a.Data.Database.Driver == b.Data.Database.Driver &&
		a.Data.Database.Connection == b.Data.Database.Connection &&
		a.Redis.Addr == b.Redis.Addr &&
		a.Redis.Password == b.Redis.Password &&
		a.Redis.DB == b.Redis.DB &&
		a.VISUAL_AI.APIKey == b.VISUAL_AI.APIKey &&
		a.VISUAL_AI.Endpoint == b.VISUAL_AI.Endpoint &&
		a.VISUAL_AI.Model == b.VISUAL_AI.Model &&
		a.VISUAL_AI.Provider == b.VISUAL_AI.Provider &&
		a.EMAIL_OTP.SMTP.Host == b.EMAIL_OTP.SMTP.Host &&
		a.EMAIL_OTP.SMTP.Port == b.EMAIL_OTP.SMTP.Port &&
		a.EMAIL_OTP.SMTP.Username == b.EMAIL_OTP.SMTP.Username &&
		a.EMAIL_OTP.SMTP.Password == b.EMAIL_OTP.SMTP.Password &&
		a.JWT.Secret == b.JWT.Secret &&
		a.JWT.AccessTTLMinutes == b.JWT.AccessTTLMinutes &&
		a.JWT.RefreshTTLHours == b.JWT.RefreshTTLHours &&
		a.Publish.Email.Enabled == b.Publish.Email.Enabled
}
