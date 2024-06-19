package config

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var c *Config

type Config struct {
	*viper.Viper
}

const (
	defaultEnvPrefix = ""
)

func init() {
	c = New(
		WithDefaultEnvVars(defaultEnvPrefix),
		WithDefaultConfigFile("", ""),
	)
}

func New(options ...Option) *Config {
	config := &Config{
		viper.New(),
	}

	for _, option := range options {
		option(config)
	}

	return config
}

func Init(config *Config) {
	c = config
}

func (cfg *Config) GetFloat32Slice(key string) []float32 {
	temp := cast.ToStringSlice(cfg.Get(key))
	ans := make([]float32, 0)
	for _, v := range temp {
		value, err := cast.ToFloat32E(v)
		if err != nil {
			return []float32{}
		}
		ans = append(ans, value)
	}
	return ans
}

func (cfg *Config) GetFloat64Slice(key string) []float64 {
	temp := cast.ToStringSlice(cfg.Get(key))
	ans := make([]float64, 0)
	for _, v := range temp {
		value, err := cast.ToFloat64E(v)
		if err != nil {
			return []float64{}
		}
		ans = append(ans, value)
	}
	return ans
}

// GetIntSlice overrides viper.GetIntSlice to cover env setup
func (cfg *Config) GetIntSlice(key string) []int {
	temp := cfg.GetStringSlice(key)
	var ans []int
	for _, v := range temp {
		value, err := cast.ToIntE(v)
		if err != nil {
			return []int{}
		}
		ans = append(ans, value)
	}

	return ans
}

func GetString(key string) string { return c.GetString(key) }

func GetBool(key string) bool { return c.GetBool(key) }

func GetInt(key string) int { return c.GetInt(key) }

func GetInt32(key string) int32 { return c.GetInt32(key) }

func GetInt64(key string) int64 { return c.GetInt64(key) }

func GetUint(key string) uint { return c.GetUint(key) }

func GetUint32(key string) uint32 { return c.GetUint32(key) }

func GetUint64(key string) uint64 { return c.GetUint64(key) }

func GetFloat64(key string) float64 { return c.GetFloat64(key) }

func GetTime(key string) time.Time { return c.GetTime(key) }

func GetDuration(key string) time.Duration { return c.GetDuration(key) }

func GetIntSlice(key string) []int { return c.GetIntSlice(key) }

func GetStringSlice(key string) []string { return c.GetStringSlice(key) }

func GetStringMap(key string) map[string]interface{} { return c.GetStringMap(key) }

func GetStringMapString(key string) map[string]string { return c.GetStringMapString(key) }

func GetStringMapStringSlice(key string) map[string][]string { return c.GetStringMapStringSlice(key) }

func GetFloat32Slice(key string) []float32 { return c.GetFloat32Slice(key) }

func GetFloat32SliceWithDefaultValue(key string, defaultValue []float32) []float32 {
	if c.IsSet(key) {
		return GetFloat32Slice(key)
	}
	return defaultValue
}

func GetFloat64Slice(key string) []float64 { return c.GetFloat64Slice(key) }

func GetFloat64SliceWithDefaultValue(key string, defaultValue []float64) []float64 {
	if c.IsSet(key) {
		return GetFloat64Slice(key)
	}
	return defaultValue
}

func GetStringWithDefaultValue(key string, defaultValue string) string {
	if c.IsSet(key) {
		return GetString(key)
	}
	return defaultValue
}

func GetBoolWithDefaultValue(key string, defaultValue bool) bool {
	if c.IsSet(key) {
		return GetBool(key)
	}
	return defaultValue
}

func GetIntWithDefaultValue(key string, defaultValue int) int {
	if c.IsSet(key) {
		return GetInt(key)
	}
	return defaultValue
}

func GetInt32WithDefaultValue(key string, defaultValue int32) int32 {
	if c.IsSet(key) {
		return GetInt32(key)
	}
	return defaultValue
}

func GetInt64WithDefaultValue(key string, defaultValue int64) int64 {
	if c.IsSet(key) {
		return GetInt64(key)
	}
	return defaultValue
}

func GetUintWithDefaultValue(key string, defaultValue uint) uint {
	if c.IsSet(key) {
		return GetUint(key)
	}
	return defaultValue
}

func GetUint32WithDefaultValue(key string, defaultValue uint32) uint32 {
	if c.IsSet(key) {
		return GetUint32(key)
	}
	return defaultValue
}

func GetUint64WithDefaultValue(key string, defaultValue uint64) uint64 {
	if c.IsSet(key) {
		return GetUint64(key)
	}
	return defaultValue
}

func GetFloat64WithDefaultValue(key string, defaultValue float64) float64 {
	if c.IsSet(key) {
		return GetFloat64(key)
	}
	return defaultValue
}

func GetTimeWithDefaultValue(key string, defaultValue time.Time) time.Time {
	if c.IsSet(key) {
		return GetTime(key)
	}
	return defaultValue
}

func GetDurationWithDefaultValue(key string, defaultValue time.Duration) time.Duration {
	if c.IsSet(key) {
		return GetDuration(key)
	}
	return defaultValue
}

func GetIntSliceWithDefaultValue(key string, defaultValue []int) []int {
	if c.IsSet(key) {
		return GetIntSlice(key)
	}
	return defaultValue
}

func GetStringSliceWithDefaultValue(key string, defaultValue []string) []string {
	if c.IsSet(key) {
		return GetStringSlice(key)
	}
	return defaultValue
}

func GetStringMapWithDefaultValue(key string, defaultValue map[string]interface{}) map[string]interface{} {
	if c.IsSet(key) {
		return GetStringMap(key)
	}
	return defaultValue
}

func GetStringMapStringWithDefaultValue(key string, defaultValue map[string]string) map[string]string {
	if c.IsSet(key) {
		return GetStringMapString(key)
	}
	return defaultValue
}

func GetStringMapStringSliceWithDefaultValue(key string, defaultValue map[string][]string) map[string][]string {
	if c.IsSet(key) {
		return GetStringMapStringSlice(key)
	}
	return defaultValue
}

func UnmarshalKey(key string, rawVal interface{}) {
	_ = c.UnmarshalKey(key, rawVal, viper.DecodeHook(jsonStringToStruct(rawVal)))
}

func jsonStringToStruct(m interface{}) func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
	return func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
		if rf != reflect.String || (rt != reflect.Struct && rt != reflect.Slice && rt != reflect.Map) {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return m, nil
		}

		err := json.Unmarshal([]byte(raw), &m)

		return m, err
	}
}

// UnmarshalToStruct unmarshal config to a config struct
func UnmarshalToStruct(key string, configStruct interface{}) error {
	bindEnvs(key, configStruct)
	return c.Unmarshal(configStruct)
}

func bindEnvs(key string, rawVal interface{}) {
	for _, k := range allKeys(key, rawVal) {
		val := c.Viper.Get(k)
		c.Viper.Set(k, val)
	}
}

func allKeys(key string, rawVal interface{}) []string {
	var b []byte
	var err error
	if key == "" {
		b, err = yaml.Marshal(rawVal)
	} else {
		b, err = yaml.Marshal(map[string]interface{}{
			key: rawVal,
		})
	}

	if err != nil {
		return nil
	}

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(b)); err != nil {
		return nil
	}

	return v.AllKeys()
}
