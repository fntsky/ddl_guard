package conf

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
)

// applyEnvVars 遍历配置结构体，用环境变量覆盖字段值
func applyEnvVars(cfg *Config) {
	if cfg == nil {
		return
	}
	applyEnvToStruct(reflect.ValueOf(cfg).Elem())
}

// applyEnvToStruct 递归处理结构体字段
func applyEnvToStruct(v reflect.Value) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 跳过不可设置的字段
		if !field.CanSet() {
			continue
		}

		// 获取 env 标签
		envName := fieldType.Tag.Get("env")

		// 如果是结构体，递归处理
		if field.Kind() == reflect.Struct {
			applyEnvToStruct(field)
			continue
		}

		// 如果没有 env 标签，跳过
		if envName == "" {
			continue
		}

		// 从环境变量读取
		envValue := os.Getenv(envName)
		if envValue == "" {
			continue
		}

		// 设置字段值
		if err := setFieldValue(field, envValue); err != nil {
			slog.Warn("failed to apply env var",
				"env", envName,
				"value", envValue,
				"error", err,
			)
		}
	}
}

// setFieldValue 将字符串值设置到字段
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("parse int: %w", err)
		}
		field.SetInt(val)
		return nil
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("parse bool: %w", err)
		}
		field.SetBool(val)
		return nil
	default:
		return fmt.Errorf("unsupported type: %s", field.Kind())
	}
}
