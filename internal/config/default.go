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
)
