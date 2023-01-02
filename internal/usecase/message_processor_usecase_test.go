package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/atjhoendz/notpushcation-service/internal/model"

	"github.com/atjhoendz/notpushcation-service/internal/model/mock"
	"github.com/golang/mock/gomock"
)

func TestMessageProcessorUsecase_ProcessNotificationMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotificationComposer := mock.NewMockComposerUsecase(ctrl)
	mockOnesignalClient := mock.NewMockOnesignalClient(ctrl)

	uc := messageProcessorUsecase{
		composers:       Composers{NotificationComposer: mockNotificationComposer},
		onesignalClient: mockOnesignalClient,
	}

	ctx := context.TODO()

	input := model.PushNotificationInput{
		Title:   "ini judul notifikasi",
		Content: "ini konten notifikasi",
		Subject: &model.Subject{Username: "joko"},
	}

	onesignalPayload := &model.OnesignalPayload{
		Contents: map[string]string{"en": "ini konten notifikasi"},
		AppID:    "112",
	}

	t.Run("success", func(t *testing.T) {
		mockNotificationComposer.EXPECT().Compose(ctx, input).Times(1).Return(onesignalPayload, nil)
		mockOnesignalClient.EXPECT().Deliver(ctx, onesignalPayload).Times(1).Return(nil)

		err := uc.ProcessNotificationMessage(ctx, input)
		require.NoError(t, err)
		assert.Nil(t, err)
	})

	t.Run("failed - onesignal error", func(t *testing.T) {
		mockNotificationComposer.EXPECT().Compose(ctx, input).Times(1).Return(onesignalPayload, nil)
		mockOnesignalClient.EXPECT().Deliver(ctx, onesignalPayload).Times(1).Return(errors.New("error onesignal"))

		err := uc.ProcessNotificationMessage(ctx, input)
		require.Error(t, err)
		assert.Equal(t, errors.New("error onesignal"), err)
	})
}
