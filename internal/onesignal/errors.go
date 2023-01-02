package onesignal

import "errors"

// errors ..
var (
	// ErrOnesignalTimeout error timeout for waiting response from onesignal
	ErrOnesignalTimeout = errors.New("onesignal_err: timeout waiting response from onesignal")
	// ErrOnesignalUnexpectedResponse error unexpected response from onesignal
	ErrOnesignalUnexpectedResponse = errors.New("onesignal_err: unexpected response from onesignal")
	// ErrOnesignalServerSide server side error from onesignal
	ErrOnesignalServerSide = errors.New("onesignal_err: onesignal server side error")
)
