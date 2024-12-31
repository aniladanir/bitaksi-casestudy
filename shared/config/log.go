package config

import "github.com/spf13/viper"

func GetLogLevel() string {
	return viper.GetString("log.level")
}

func GetAccessLogFile() string {
	// /var/log/matching-api/access.log
	return viper.GetString("log.access.file")
}

func GetLogFile() string {
	// /var/log/matching-api/matching-api.log
	return viper.GetString("log.file")
}

func IsDebug() bool {
	return viper.GetString("log.level") == "debug"
}

func GetLogMaxAgeInDays() int {
	return viper.GetInt("log.maxAge")
}

func GetLogMaxSizeInMB() int {
	return viper.GetInt("log.maxSize")
}

func GetLogMaxBackups() int {
	return viper.GetInt("log.maxBackups")
}

func GetLogGzipArchive() bool {
	return viper.GetBool("log.gzipArchive")
}
