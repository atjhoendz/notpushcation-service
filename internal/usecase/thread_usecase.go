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
	//sseBroker  *model.SSEBroker
	sseServer *sse.Server
}

func NewThreadUsecase(threadRepo model.ThreadRepository, sseServer *sse.Server) model.ThreadUsecase {
	return &threadUsecase{threadRepo: threadRepo, sseServer: sseServer}
}

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

	//msg := &model.SSEMessage{
	//	Event: model.Update,
	//	Data:  res.ToJSONString(),
	//}
	//
	//u.sseBroker.BroadcastMessage(msg)

	msg := &sse.Event{
		Data:  res.ToArrayByte(),
		Event: []byte(model.Update),
	}

	u.sseServer.Publish("thread", msg)

	return res, nil
}

func (u threadUsecase) FindAll(ctx context.Context) (threads []*model.Thread, err error) {
	//TODO implement me
	panic("implement me")
}
