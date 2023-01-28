package config

import "time"

// default constants
const (
	DefaultTLSHandshakeTimeout   = 5 * time.Second
	DefaultTLSInsecureSkipVerify = true
	DefaultHTTPTimeout           = 1 * time.Second

	DefaultCockroachMaxIdleConns    = 3
	DefaultCockroachMaxOpenConns    = 5
	DefaultCockroachConnMaxLifetime = 1 * time.Hour
	DefaultCockroachPingInterval    = 1 * time.Second
	DefaultCockroachRetryAttempts   = 3

	NATSDurableID                       = "sse-service-durable"
	DefaultNATSJSRetryOnFailedConnect   = true
	DefaultNATSJSMaxReconnect           = -1
	DefaultNATSJSReconnectWait          = 1 * time.Second
	DefaultNATSJSRetryAttempts          = 3
	DefaultNATSJSRetryInterval          = 2 * time.Second
	DefaultNATSJSSubscribeRetryAttempts = 3
	DefaultNATSJSSubscribeRetryInterval = 2 * time.Second
	DefaultNATSJSStreamMaxAge           = 1 * 24 * time.Hour
	DefaultNATSJSStreamMaxMessages      = 100000

	DefaultRateLimiterPeriod       = time.Second
	DefaultRateLimiterRequestLimit = 20

	DefaultRedisIdleTimeout     = 240 * time.Second
	DefaultRedisMaxConnLifetime = time.Minute
)
