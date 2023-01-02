package config

import (
	"github.com/kumparan/go-connect"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"strings"
)

// Env :nodoc:
func Env() string {
	return viper.GetString("env")
}

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}

// HTTPPort :nodoc:
func HTTPPort() string {
	return viper.GetString("ports.http")
}

// EnableIntrospection :nodoc:
func EnableIntrospection() bool {
	return viper.GetBool("enable_introspection")
}

// OnesignalAPIKey :nodoc:
func OnesignalAPIKey() string {
	return viper.GetString("onesignal.api_key")
}

// OnesignalAppID :nodoc:
func OnesignalAppID() string {
	return viper.GetString("onesignal.app_id")
}

// DefaultHTTPOptions :nodoc:
func DefaultHTTPOptions() *connect.HTTPConnectionOptions {
	return &connect.HTTPConnectionOptions{
		TLSHandshakeTimeout:   DefaultTLSHandshakeTimeout,
		TLSInsecureSkipVerify: DefaultTLSInsecureSkipVerify,
		Timeout:               DefaultHTTPTimeout,
	}
}

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("./../../..")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("svc")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Warningf("%v", err)
	}
}
