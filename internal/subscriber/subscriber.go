package subscriber

import (
	"context"
	"encoding/json"

	"github.com/atjhoendz/notpushcation-service/internal/model"

	"github.com/atjhoendz/notpushcation-service/event"
	"github.com/kumparan/ferstream"
)

const (
	contextCaller     = caller("sse-service-subscriber-call")
	contentQueueGroup = "sse-service-content-queue"
)

type (
	caller string

	// Subscriber :nodoc:
	Subscriber interface {
		RegisterNATSJetStream(js ferstream.JetStream)
		SubscribeJetStreamEvent() error
	}

	subscriberImpl struct {
		js                  ferstream.JetStream
		liveBlogPostUsecase model.LiveBlogPostUsecase
	}
)

// NewSubscriber :nodoc:
func NewSubscriber(liveBlogPostUsecase model.LiveBlogPostUsecase) Subscriber {
	return &subscriberImpl{liveBlogPostUsecase: liveBlogPostUsecase}
}

// RegisterNATSJetStream :nodoc:
func (s *subscriberImpl) RegisterNATSJetStream(js ferstream.JetStream) {
	s.js = js
}

// SubscribeJetStreamEvent :nodoc:
func (s *subscriberImpl) SubscribeJetStreamEvent() error {
	eventHandler := createJetStreamEventHandler(event.LiveBlogPostStreamName, func(ctx context.Context, msg *ferstream.NatsEventMessage) error {
		switch msg.NatsEvent.GetSubject() {
		case event.LiveBlogPostSubjectCreate:
			var body model.CreateLiveBlogPostInput
			_ = json.Unmarshal([]byte(msg.Body), &body)

			s.liveBlogPostUsecase.HandleEvent(body)
			return nil
		default:
			return nil
		}
	})

	err := subscribeStream(
		s.js,
		event.LiveBlogPostStreamName,
		event.LiveBlogPostSubjectAll,
		contentQueueGroup,
		eventHandler,
		nil)
	if err != nil {
		return err
	}

	return nil
}
