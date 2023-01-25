package model

import (
	"encoding/json"

	"github.com/kumparan/ferstream"
)

type (
	// LiveBlogPostUsecase :nodoc:
	LiveBlogPostUsecase interface {
		RegisterNATSJetStream(js ferstream.JetStream)
		Create(input CreateLiveBlogPostInput)
		HandleEvent(input CreateLiveBlogPostInput)
	}

	// CreateLiveBlogPostInput :nodoc:
	CreateLiveBlogPostInput struct {
		ThreadID int64  `json:"thread_id"`
		Title    string `json:"title"`
	}
)

// ToArrayByte :nodoc:
func (c CreateLiveBlogPostInput) ToArrayByte() []byte {
	res, _ := json.Marshal(c)
	return res
}

const (
	// LiveBlogPost :nodoc:
	LiveBlogPost SSEEvent = "LIVE_BLOG_POST"
)
