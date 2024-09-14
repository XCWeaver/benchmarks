package main

import (
	"context"
	"sync"
	"time"

	"github.com/TiagoMalhadas/xcweaver"
)

// Post_storage_america component.
type Post_storage_america interface {
	GetPost(context.Context, Post_id_obj) (string, error)
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
}

// Implementation of the Post_storage_america component.
type post_storage_america struct {
	xcweaver.Implements[Post_storage_america]
	clientRedis             xcweaver.Antipode
	mu                      sync.Mutex
	consistencyWindowValues []float64
}

type Post_id_obj struct {
	xcweaver.AutoMarshal
	PostId    string
	WriteTime int64
}

func (p *post_storage_america) Init(ctx context.Context) error {
	logger := p.Logger(ctx)

	logger.Info("post storage service at us running!")
	return nil
}

func (p *post_storage_america) GetPost(ctx context.Context, postId Post_id_obj) (string, error) {
	logger := p.Logger(ctx)

	startTimeMs := time.Now().UnixMilli()
	post, _, err := p.clientRedis.Read(ctx, "posts", postId.PostId)
	readPostDurationMs.Put(float64(time.Now().UnixMilli() - startTimeMs))
	consistencyWindowMs := float64(time.Now().UnixMilli() - postId.WriteTime)
	consistencyWindow.Put(consistencyWindowMs)
	p.mu.Lock()
	p.consistencyWindowValues = append(p.consistencyWindowValues, consistencyWindowMs)
	p.mu.Unlock()
	if err == xcweaver.ErrNotFound {
		inconsistencies.Inc()
		logger.Error("post not found")
		return "post not found", err
	} else if err != nil {
		return "", err
	}

	return post, nil
}

func (p *post_storage_america) GetConsistencyWindowValues(ctx context.Context) ([]float64, error) {
	logger := p.Logger(ctx)
	logger.Debug("entering GetConsistencyWindowValues")
	p.mu.Lock()
	values := p.consistencyWindowValues
	p.mu.Unlock()
	return values, nil
}
