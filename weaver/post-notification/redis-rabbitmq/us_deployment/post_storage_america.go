package main

import (
	"context"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/redis/go-redis/v9"
)

// Post_storage_america component.
type Post_storage_america interface {
	GetPost(context.Context, Post_id_obj) (string, error)
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
}

// Implementation of the Post_storage_america component.
type post_storage_america struct {
	weaver.Implements[Post_storage_america]
	weaver.WithConfig[post_storage_usOptions]
	client                  *redis.Client
	mu                      sync.Mutex
	consistencyWindowValues []float64
}

type Post_id_obj struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}

type post_storage_usOptions struct {
	RedisAddr     string `toml:"redis_address"`
	RedisPort     string `toml:"redis_port"`
	RedisPassword string `toml:"redis_password"`
}

func (p *post_storage_america) Init(ctx context.Context) error {
	logger := p.Logger(ctx)

	p.client = redis.NewClient(&redis.Options{
		Addr:     p.Config().RedisAddr + ":" + p.Config().RedisPort,
		Password: p.Config().RedisPassword,
		DB:       0,                        // use default DB
	})

	logger.Info("post storage service at eu running!", "redis host", p.Config().RedisAddr, "port", p.Config().RedisPort, "password", p.Config().RedisPassword)

	return nil
}

func (p *post_storage_america) GetPost(ctx context.Context, postId Post_id_obj) (string, error) {
	logger := p.Logger(ctx)

	logger.Debug("Reading post!", "postId", postId)
	startTimeMs := time.Now().UnixMilli()
	post, err := p.client.Get(ctx, postId.PostId).Result()
	readPostDurationMs.Put(float64(time.Now().UnixMilli() - startTimeMs))
	consistencyWindowMs := float64(time.Now().UnixMilli() - postId.WriteTime)
	consistencyWindow.Put(consistencyWindowMs)
	p.mu.Lock()
	p.consistencyWindowValues = append(p.consistencyWindowValues, consistencyWindowMs)
	p.mu.Unlock()
	if err == redis.Nil {
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
