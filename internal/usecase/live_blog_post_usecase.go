package usecase

import (
	"fmt"

	"github.com/atjhoendz/notpushcation-service/internal/config"

	"github.com/nats-io/nats.go"

	log "github.com/sirupsen/logrus"

	"github.com/atjhoendz/notpushcation-service/event"

	"github.com/kumparan/ferstream"

	"github.com/atjhoendz/notpushcation-service/internal/model"
	"github.com/r3labs/sse/v2"
)

type liveBlogPostUsecase struct {
	sseServer *sse.Server
	js        ferstream.JetStream
}

// NewLiveBlogPostUsecase :nodoc:
func NewLiveBlogPostUsecase(sseServer *sse.Server) model.LiveBlogPostUsecase {
	return &liveBlogPostUsecase{sseServer: sseServer}
}

// RegisterNATSJetStream :nodoc:
func (l *liveBlogPostUsecase) RegisterNATSJetStream(js ferstream.JetStream) {
	l.js = js
}

// InitStream :nodoc:
func (l *liveBlogPostUsecase) InitStream() error {
	_, err := l.js.AddStream(&nats.StreamConfig{
		Name:     event.LiveBlogPostStreamName,
		Subjects: []string{event.LiveBlogPostSubjectAll},
		MaxAge:   config.NATSJSStreamMaxAge(),
		MaxMsgs:  config.NATSJSStreamMaxMessages(),
		Storage:  nats.FileStorage,
	})
	if err != nil {
		log.Error(err)
	}

	return err
}

// Create :nodoc:
func (l *liveBlogPostUsecase) Create(input model.CreateLiveBlogPostInput) {
	msg := ferstream.NatsEventMessage{
		NatsEvent: &ferstream.NatsEvent{
			ID: input.ThreadID,
		},
	}

	msg.WithBody(input)

	msgByte, err := msg.Build()
	if err != nil {
		log.Error(err)
		return
	}

	_, err = l.js.Publish(event.LiveBlogPostSubjectCreate, msgByte)
	if err != nil {
		log.Error(err)
	}
}

// HandleEvent :nodoc:
func (l *liveBlogPostUsecase) HandleEvent(input model.CreateLiveBlogPostInput) {
	msg := &sse.Event{
		Data:  input.ToArrayByte(),
		Event: []byte(model.LiveBlogPost),
	}

	l.sseServer.Publish(fmt.Sprintf(model.StreamIDLiveBlogPostByThreadIDTempl, input.ThreadID), msg)
}
