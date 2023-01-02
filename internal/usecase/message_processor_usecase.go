package usecase

import (
	"context"

	"github.com/atjhoendz/notpushcation-service/internal/model"
	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

type (
	// Composers :nodoc:
	Composers struct {
		NotificationComposer model.ComposerUsecase
	}

	messageProcessorUsecase struct {
		composers       Composers
		onesignalClient model.OnesignalClient
	}
)

// NewMessageProcessorUsecase :nodoc:
func NewMessageProcessorUsecase(composers Composers, onesignalClient model.OnesignalClient) model.MessageProcessorUsecase {
	return &messageProcessorUsecase{
		composers:       composers,
		onesignalClient: onesignalClient,
	}
}

// ProcessNotificationMessage :nodoc:
func (m *messageProcessorUsecase) ProcessNotificationMessage(ctx context.Context, input model.PushNotificationInput) error {
	logger := log.WithFields(log.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": input,
	})

	err := m.createNotification(ctx, input)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (m *messageProcessorUsecase) createNotification(ctx context.Context, input model.PushNotificationInput) error {
	logger := log.WithFields(log.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": input,
	})

	pushNotificationMessage, err := m.composers.NotificationComposer.Compose(ctx, input)
	if err != nil {
		logger.Error(err)
		return err
	}

	if pushNotificationMessage == nil {
		return nil
	}

	err = m.onesignalClient.Deliver(ctx, pushNotificationMessage)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
