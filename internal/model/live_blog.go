package model

import (
	"encoding/json"

	"github.com/kumparan/ferstream"
)

type (
	LiveBlogPostUsecase interface {
		RegisterNATSJetStream(js ferstream.JetStream)
		Create(input CreateLiveBlogPostInput)
		HandleEvent(input CreateLiveBlogPostInput)
	}

	CreateLiveBlogPostInput struct {
		ThreadID int64  `json:"thread_id"`
		Title    string `json:"title"`
	}
)

func (c CreateLiveBlogPostInput) ToArrayByte() []byte {
	res, _ := json.Marshal(c)
	return res
}
