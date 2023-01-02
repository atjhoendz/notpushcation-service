package config

import "time"

// default constants
const (
	DefaultTLSHandshakeTimeout   = 5 * time.Second
	DefaultTLSInsecureSkipVerify = true
	DefaultHTTPTimeout           = 1 * time.Second
)
