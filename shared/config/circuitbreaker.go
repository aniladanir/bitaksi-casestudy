package config

import "github.com/spf13/viper"

func GetCircuitBreakerRetryTimeout() int {
	return viper.GetInt("circuitBreaker.retryTimeout")
}

func GetCircuitBreakerMaxFailures() int {
	return viper.GetInt("circuitBreaker.maxFailures")
}
