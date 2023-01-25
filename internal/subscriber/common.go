package subscriber

import (
	"context"

	"github.com/atjhoendz/notpushcation-service/internal/config"
	"github.com/kumparan/ferstream"
	"github.com/kumparan/go-utils"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func getMsgFromFerstream(payload ferstream.MessageParser) (msg *ferstream.NatsEventMessage, err error) {
	logger := log.WithFields(log.Fields{
		"payload": utils.Dump(payload),
	})

	msg, ok := payload.(*ferstream.NatsEventMessage)
	if !ok {
		err = ferstream.ErrCastingPayloadToStruct
		logger.Error(err)
		return
	}
	return
}

func createJetStreamEventHandler(streamName string, fn func(ctx context.Context, msg *ferstream.NatsEventMessage) error) ferstream.MessageHandler {
	return func(payload ferstream.MessageParser) (err error) {
		logger := log.WithFields(log.Fields{
			"payload": utils.Dump(payload),
		})
		ctx := context.WithValue(context.Background(), contextCaller, streamName)

		msg, err := getMsgFromFerstream(payload)
		if err != nil {
			logger.Error(err)
			return err
		}

		return fn(ctx, msg)
	}
}

func subscribeStream(js ferstream.JetStream, streamName, streamSubject, queueGroup string, eventHandler, eventErrorHandler ferstream.MessageHandler) error {
	logger := log.WithFields(log.Fields{
		"streamName":    streamName,
		"streamSubject": streamSubject,
		"queueGroup":    queueGroup,
	})

	err := utils.Retry(config.NATSJSSubscribeRetryAttempts(), config.NATSJSSubscribeRetryInterval(), func() error {
		_, err := js.QueueSubscribe(streamSubject, queueGroup,
			ferstream.NewNATSMessageHandler(
				new(ferstream.NatsEventMessage),
				config.NATSJSRetryAttempts(),
				config.NATSJSRetryInterval(),
				eventHandler,
				eventErrorHandler),
			nats.ManualAck(),
			nats.Durable(config.NATSDurableID),
		)
		return err
	})
	if err != nil {
		logger.Error(err)
	}

	return err
}
