package main

import (
	"context"
	"time"

	"github.com/TiagoMalhadas/xcweaver"
	"github.com/google/uuid"
)

// Post_storage component.
type Post_storage_europe interface {
	Post(context.Context, string) ([]byte, Post_id_obj, error)
}

// Implementation of the Post_storage component.
type post_storage_europe struct {
	xcweaver.Implements[Post_storage_europe]
	clientRedis xcweaver.Antipode
}

type Post_id_obj struct {
	xcweaver.AutoMarshal
	PostId    string
	WriteTime int64
}

func (p *post_storage_europe) Init(ctx context.Context) error {
	logger := p.Logger(ctx)
	logger.Info("post storage service at eu running!")

	return nil
}

func (p *post_storage_europe) Post(ctx context.Context, post string) ([]byte, Post_id_obj, error) {
	logger := p.Logger(ctx)

	id := uuid.New()

	writeStartTimeMs := time.Now().UnixMilli()
	ctx, err := p.clientRedis.Write(ctx, "posts", id.String(), post)
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
