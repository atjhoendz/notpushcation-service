package usecase

import (
	"context"

	"github.com/atjhoendz/notpushcation-service/internal/model"
)

type notificationComposer struct{}

// NewNotificationComposer :nodoc:
func NewNotificationComposer() model.ComposerUsecase {
	return &notificationComposer{}
}

func (c *notificationComposer) Compose(ctx context.Context, input model.PushNotificationInput) (*model.OnesignalPayload, error) {
	var segments []string
	for _, val := range input.Segments {
		segments = append(segments, val.GetString())
	}

	return &model.OnesignalPayload{
		Headings:         map[string]string{"en": input.Title},
		Contents:         map[string]string{"en": input.Content},
		IncludedSegments: segments,
	}, nil
}
