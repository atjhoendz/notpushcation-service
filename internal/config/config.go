package config

import (
	"fmt"
	"time"

	"github.com/kumparan/go-utils"

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

// CockroachHost :nodoc:
func CockroachHost() string {
	return viper.GetString("cockroach.host")
}

// CockroachDatabase :nodoc:
func CockroachDatabase() string {
	return viper.GetString("cockroach.database")
}

// CockroachUsername :nodoc:
func CockroachUsername() string {
	return viper.GetString("cockroach.username")
}

// CockroachPassword :nodoc:
func CockroachPassword() string {
	return viper.GetString("cockroach.password")
}

// CockroachSSLMode :nodoc:
func CockroachSSLMode() string {
	if viper.IsSet("cockroach.sslmode") {
		return viper.GetString("cockroach.sslmode")
	}
	return "disable"
}

// CockroachMaxIdleConns :nodoc:
func CockroachMaxIdleConns() int {
	if viper.GetInt("cockroach.max_idle_conns") <= 0 {
		return DefaultCockroachMaxIdleConns
	}
	return viper.GetInt("cockroach.max_idle_conns")
}

// CockroachMaxOpenConns :nodoc:
func CockroachMaxOpenConns() int {
	if viper.GetInt("cockroach.max_open_conns") <= 0 {
		return DefaultCockroachMaxOpenConns
	}
	return viper.GetInt("cockroach.max_open_conns")
}

// CockroachConnMaxLifetime :nodoc:
func CockroachConnMaxLifetime() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("cockroach.conn_max_lifetime"), DefaultCockroachConnMaxLifetime)
}

// CockroachPingInterval :nodoc:
func CockroachPingInterval() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("cockroach.ping_interval"), DefaultCockroachPingInterval)
}

// CockroachRetryAttempts :nodoc:
func CockroachRetryAttempts() int {
	if viper.GetInt("cockroach.retry_attempts") > 0 {
		return viper.GetInt("cockroach.retry_attempts")
	}
	return DefaultCockroachRetryAttempts
}

// DatabaseDSN :nodoc:
func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		CockroachUsername(),
		CockroachPassword(),
		CockroachHost(),
		CockroachDatabase(),
		CockroachSSLMode())
}

// NATSJSHost :nodoc:
func NATSJSHost() string {
	return viper.GetString("nats_js.host")
}

// NATSJSRetryOnFailedConnect :nodoc:
func NATSJSRetryOnFailedConnect() bool {
	if !viper.IsSet("nats_js.retry_on_failed_connect") {
		return DefaultNATSJSRetryOnFailedConnect
	}
	return viper.GetBool("nats_js.retry_on_failed_connect")
}

// NATSJSMaxReconnect :nodoc:
func NATSJSMaxReconnect() int {
	return utils.ValueOrDefault[int](viper.GetInt("nats_js.max_reconnect"), DefaultNATSJSMaxReconnect)
}

// NATSJSReconnectWait :nodoc:
func NATSJSReconnectWait() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("nats_js.reconnect_wait"), DefaultNATSJSReconnectWait)
}

// NATSJSRetryAttempts :nodoc:
func NATSJSRetryAttempts() int {
	return utils.ValueOrDefault[int](viper.GetInt("nats_js.retry_attempts"), DefaultNATSJSRetryAttempts)
}

// NATSJSRetryInterval :nodoc:
func NATSJSRetryInterval() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("nats_js.retry_interval"), DefaultNATSJSRetryInterval)
}

// NATSJSSubscribeRetryAttempts :nodoc:
func NATSJSSubscribeRetryAttempts() int {
	return utils.ValueOrDefault[int](viper.GetInt("nats_js.subscribe_retry_attempts"), DefaultNATSJSSubscribeRetryAttempts)
}

// NATSJSSubscribeRetryInterval :nodoc:
func NATSJSSubscribeRetryInterval() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("nats_js.subscribe_retry_interval"), DefaultNATSJSSubscribeRetryInterval)
}

// NATSJSStreamMaxAge :nodoc:
func NATSJSStreamMaxAge() time.Duration {
	return utils.ParseDurationWithDefault(viper.GetString("nats_js.stream_max_age"), DefaultNATSJSStreamMaxAge)
}

// NATSJSStreamMaxMessages :nodoc:
func NATSJSStreamMaxMessages() int64 {
	return utils.ValueOrDefault[int64](viper.GetInt64("nats_js.stream_max_messages"), DefaultNATSJSStreamMaxMessages)
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
