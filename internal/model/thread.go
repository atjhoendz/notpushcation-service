package model

import (
	"context"
	"encoding/json"
)

type (
	// ThreadRepository :nodoc:
	ThreadRepository interface {
		Create(ctx context.Context, t *Thread) (*Thread, error)
	}

	// ThreadUsecase :nodoc:
	ThreadUsecase interface {
		Create(ctx context.Context, input CreateThreadInput) (*Thread, error)
	}

	// Thread :nodoc:
	Thread struct {
		ID      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	// CreateThreadInput :nodoc:
	CreateThreadInput struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	// SSEEvent :nodoc:
	SSEEvent string
)

// ToArrayByte :nodoc:
func (t Thread) ToArrayByte() []byte {
	res, _ := json.Marshal(t)
	return res
}

const (
	// Update :nodoc:
	Update SSEEvent = "UPDATE"
)
