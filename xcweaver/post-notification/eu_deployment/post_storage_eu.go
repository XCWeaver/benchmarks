package main

import (
	"context"
	"time"

	"github.com/XCWeaver/xcweaver"
	"github.com/google/uuid"
)

// PostStorage component.
type PostStorageEu interface {
	Post(context.Context, string) ([]byte, Post_id_obj, error)
}

// Implementation of the Post_storage component.
type postStorageEu struct {
	xcweaver.Implements[PostStorageEu]
	client xcweaver.Antipode
}

type Post_id_obj struct {
	xcweaver.AutoMarshal
	PostId    string
	WriteTime int64
}

func (p *postStorageEu) Init(ctx context.Context) error {
	logger := p.Logger(ctx)
	logger.Info("post storage service at eu running!")

	return nil
}

func (p *postStorageEu) Post(ctx context.Context, post string) ([]byte, Post_id_obj, error) {
	logger := p.Logger(ctx)

	id := uuid.New()

	writeStartTimeMs := time.Now().UnixMilli()
	ctx, err := p.client.Write(ctx, "posts", id.String(), post)
	writePostDurationMs.Put(float64(time.Now().UnixMilli() - writeStartTimeMs))
	if err != nil {
		logger.Error("Error writing post!", "msg", err.Error())
		return []byte{}, Post_id_obj{}, err
	}
	logger.Debug("Post successfully stored!", "postId", id.String(), "post", post)

	lineage, err := xcweaver.GetLineage(ctx)
	if err != nil {
		logger.Error("Error getting lineage from context!", "msg", err.Error())
		return []byte{}, Post_id_obj{}, err
	}

	return lineage, Post_id_obj{PostId: id.String(), WriteTime: writeStartTimeMs}, nil
}
