package model

import "context"

// MessageProcessorUsecase :nodoc:
type MessageProcessorUsecase interface {
	ProcessNotificationMessage(ctx context.Context, input PushNotificationInput) error
}

type (
	// PushNotificationInput :nodoc:
	PushNotificationInput struct {
		Title    string             `json:"title"`
		Content  string             `json:"content"`
		Segments []OnesignalSegment `json:"segments"`
	}
)
