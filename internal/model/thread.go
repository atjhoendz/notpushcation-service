package model

import (
	"bytes"
	"context"
	"encoding/json"
)

type (
	ThreadRepository interface {
		Create(ctx context.Context, t *Thread) (*Thread, error)
		FindAll(ctx context.Context) (threads []*Thread, err error)
	}

	ThreadUsecase interface {
		Create(ctx context.Context, input CreateThreadInput) (*Thread, error)
		FindAll(ctx context.Context) (threads []*Thread, err error)
	}

	Thread struct {
		ID      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	CreateThreadInput struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
)

func (t Thread) ToJSONString() string {
	var buf bytes.Buffer
	res, _ := json.Marshal(t)
	buf.Write(res)
	return buf.String()
}

func (t Thread) ToArrayByte() []byte {
	res, _ := json.Marshal(t)
	return res
}
