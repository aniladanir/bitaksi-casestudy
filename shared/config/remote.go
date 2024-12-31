package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetRemoteUrl(serviceName string) string {
	return viper.GetString(fmt.Sprintf("remote.%s.url", serviceName))
}

func GetRemoteVersion(serviceName string) string {
	return viper.GetString(fmt.Sprintf("remote.%s.version", serviceName))
}
