package model

import "context"

// MessageProcessorUsecase :nodoc:
type MessageProcessorUsecase interface {
	ProcessNotificationMessage(ctx context.Context, input PushNotificationInput) error
}

// Subject :nodoc:
type (
	Subject struct {
		Username string
	}

	// PushNotificationInput :nodoc:
	PushNotificationInput struct {
		Title    string             `json:"title"`
		Content  string             `json:"content"`
		Subject  *Subject           `json:"subject"`
		Segments []OnesignalSegment `json:"segments"`
	}
)
