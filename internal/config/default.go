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

	NATSDurableID = "sse-service-durable"
	// DefaultNATSJSRetryOnFailedConnect :nodoc:
	DefaultNATSJSRetryOnFailedConnect = true
	// DefaultNATSJSMaxReconnect :nodoc:
	DefaultNATSJSMaxReconnect = -1
	// DefaultNATSJSReconnectWait :nodoc:
	DefaultNATSJSReconnectWait = 1 * time.Second
	// DefaultNATSJSRetryAttempts :nodoc:
	DefaultNATSJSRetryAttempts = 3
	// DefaultNATSJSRetryInterval :nodoc:
	DefaultNATSJSRetryInterval = 2 * time.Second
	// DefaultNATSJSSubscribeRetryAttempts :nodoc:
	DefaultNATSJSSubscribeRetryAttempts = 3
	// DefaultNATSJSSubscribeRetryInterval :nodoc:
	DefaultNATSJSSubscribeRetryInterval = 2 * time.Second
	// DefaultNATSJSStreamMaxAge :nodoc:
	DefaultNATSJSStreamMaxAge = 1 * 24 * time.Hour
	// DefaultNATSJSDeliveryTimeInMinute :nodoc:
	DefaultNATSJSDeliveryTimeInMinute = -30
	// DefaultNATSJSStreamMaxMessages :nodoc:
	DefaultNATSJSStreamMaxMessages = 100000
)
