package config

import (
	"github.com/spf13/viper"
)

func SetDefaults(defaults map[string]string) {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
}

func Init(file string) error {
	viper.SetConfigFile(file)
	viper.WatchConfig()

	return viper.ReadInConfig()
}

func Get(key string) any {
	return viper.Get(key)
}

func GetAPIVersion() string {
	return viper.GetString("api.version")
}
