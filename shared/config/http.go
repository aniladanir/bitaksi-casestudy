package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetHttpServerAddress() string {
	return fmt.Sprintf("%s:%d",
		viper.GetString("http.ipAddress"),
		viper.GetInt("http.port"),
	)
}

func GetHttpReadTimeout() int {
	return viper.GetInt("http.readTimeout")
}

func GetHttpWriteTimeout() int {
	return viper.GetInt("http.writeTimeout")
}

func GetHttpIdleTimeout() int {
	return viper.GetInt("http.idleTimeout")
}

func GetHttpClientTimeout() int {
	return viper.GetInt("http.clientTimeout")
}
