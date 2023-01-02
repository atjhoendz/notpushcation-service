package model

import (
	"context"
)

// ComposerUsecase :nodoc:
type ComposerUsecase interface {
	Compose(ctx context.Context, input PushNotificationInput) (*OnesignalPayload, error)
}
