package config

import "github.com/spf13/viper"

func GetDBName() string {
	return viper.GetString("db.name")
}

func GetDBConnectionString() string {
	return viper.GetString("db.connectionString")
}
