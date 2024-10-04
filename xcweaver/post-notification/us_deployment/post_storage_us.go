package main

import (
	"context"
	"sync"
	"time"

	"github.com/XCWeaver/xcweaver"
)

// PostStorageUS component.
type PostStorageUS interface {
	GetPost(context.Context, Post_id_obj) (string, error)
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
	GetInconsistencies(ctx context.Context) (int, error)
	Reset(ctx context.Context) error
}

// Implementation of the PostStorageUS component.
type postStorageUS struct {
	xcweaver.Implements[PostStorageUS]
	client                  xcweaver.Antipode
	mu                      sync.Mutex
	consistencyWindowValues []float64
	muInc                   sync.Mutex
	inconsistencies         int
}

type Post_id_obj struct {
	xcweaver.AutoMarshal
	PostId    string
	WriteTime int64
}

func (p *postStorageUS) Init(ctx context.Context) error {
	logger := p.Logger(ctx)

	p.inconsistencies = 0

	logger.Info("post storage service at us running!")
	return nil
}

func (p *postStorageUS) GetPost(ctx context.Context, postId Post_id_obj) (string, error) {
	logger := p.Logger(ctx)

	startTimeMs := time.Now().UnixMilli()
	post, _, err := p.client.Read(ctx, "posts", postId.PostId)
	readPostDurationMs.Put(float64(time.Now().UnixMilli() - startTimeMs))
	consistencyWindowMs := float64(time.Now().UnixMilli() - postId.WriteTime)
	consistencyWindow.Put(consistencyWindowMs)
	p.mu.Lock()
	p.consistencyWindowValues = append(p.consistencyWindowValues, consistencyWindowMs)
	p.mu.Unlock()
	if err == xcweaver.ErrNotFound {
		inconsistencies.Inc()
		p.muInc.Lock()
		p.inconsistencies += 1
		p.muInc.Unlock()
		logger.Error("post not found")
		return "post not found", err
	} else if err != nil {
		return "", err
	}

	return post, nil
}

func (p *postStorageUS) GetConsistencyWindowValues(ctx context.Context) ([]float64, error) {
	logger := p.Logger(ctx)
	logger.Debug("entering GetConsistencyWindowValues")
	p.mu.Lock()
	values := p.consistencyWindowValues
	p.mu.Unlock()
	return values, nil
}

func (p *postStorageUS) GetInconsistencies(ctx context.Context) (int, error) {
	logger := p.Logger(ctx)
	logger.Debug("entering GetInconsistencies")
	p.muInc.Lock()
	inconsistencies := p.inconsistencies
	p.muInc.Unlock()
	return inconsistencies, nil
}

func (p *postStorageUS) Reset(ctx context.Context) error {
	logger := p.Logger(ctx)
	logger.Debug("entering Reset")
	p.inconsistencies = 0
	p.consistencyWindowValues = []float64{}
	return nil
}
