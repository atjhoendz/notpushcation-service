package model

import (
	"context"
)

type (
	// OnesignalClient :nodoc:
	OnesignalClient interface {
		Deliver(ctx context.Context, message *OnesignalPayload) error
	}

	// DataPayload :nodoc:
	//DataPayload struct {
	//	Username string `json:"username,omitempty"`
	//}

	// OnesignalPayload represents a onesignal notification payload
	OnesignalPayload struct {
		Headings         map[string]string `json:"headings"`
		Contents         map[string]string `json:"contents"`
		AppID            string            `json:"app_id"`
		IncludedSegments []string          `json:"included_segments"`
		//Data     DataPayload       `json:"data"`
	}

	// OnesignalSegment :nodoc:
	OnesignalSegment string
)

const (
	// SubscribedUsers :nodoc:
	SubscribedUsers OnesignalSegment = "SUBSCRIBED_USERS"
	// ActiveUsers :nodoc:
	ActiveUsers OnesignalSegment = "ACTIVE_USERS"
	// InactiveUsers :nodoc:
	InactiveUsers OnesignalSegment = "INACTIVE_USERS"
)

// GetString :nodoc:
func (s OnesignalSegment) GetString() string {
	mapSegment := map[OnesignalSegment]string{
		SubscribedUsers: "Subscribed Users",
		ActiveUsers:     "Active Users",
		InactiveUsers:   "Inactive Users",
	}

	return mapSegment[s]
}
