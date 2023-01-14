package usecase

import (
	"github.com/atjhoendz/notpushcation-service/internal/model"
)

type SSEBrokerUsecase struct {
	broker *model.SSEBroker
}

func NewSSEBrokerUsecase(b *model.SSEBroker) model.SSEBrokerUsecase {
	return &SSEBrokerUsecase{
		broker: b,
	}
}

func (s *SSEBrokerUsecase) Serve() {}
