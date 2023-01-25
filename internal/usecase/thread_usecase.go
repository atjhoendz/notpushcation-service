package usecase

import (
	"context"

	"github.com/r3labs/sse/v2"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"

	"github.com/atjhoendz/notpushcation-service/internal/model"
)

type threadUsecase struct {
	threadRepo model.ThreadRepository
	sseServer  *sse.Server
}

// NewThreadUsecase :nodoc:
func NewThreadUsecase(threadRepo model.ThreadRepository, sseServer *sse.Server) model.ThreadUsecase {
	return &threadUsecase{threadRepo: threadRepo, sseServer: sseServer}
}

// Create :nodoc:
func (u threadUsecase) Create(ctx context.Context, input model.CreateThreadInput) (*model.Thread, error) {
	logger := log.WithFields(log.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	newThread := &model.Thread{
		Title:   input.Title,
		Content: input.Content,
	}

	res, err := u.threadRepo.Create(ctx, newThread)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	msg := &sse.Event{
		Data:  res.ToArrayByte(),
		Event: []byte(model.Update),
	}

	u.sseServer.Publish("thread", msg)

	return res, nil
}
