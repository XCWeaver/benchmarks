package main

import (
	"context"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// PostStorageEu component.
type PostStorageEu interface {
	Post(context.Context, string) (Post_id_obj, error)
}

type Post_id_obj struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}

// Implementation of the PostStorageEu component.
type postStorageEu struct {
	weaver.Implements[PostStorageEu]
	weaver.WithConfig[postStorageEuOptions]
	client *redis.Client
}

type postStorageEuOptions struct {
	RedisAddr     string `toml:"redis_address"`
	RedisPort     string `toml:"redis_port"`
	RedisPassword string `toml:"redis_password"`
}

func (p *postStorageEu) Init(ctx context.Context) error {
	logger := p.Logger(ctx)

	p.client = redis.NewClient(&redis.Options{
		Addr:     p.Config().RedisAddr + ":" + p.Config().RedisPort,
		Password: p.Config().RedisPassword,
		DB:       0, // use default DB
	})

	logger.Info("post storage service at eu running!", "redis host", p.Config().RedisAddr, "port", p.Config().RedisPort, "password", p.Config().RedisPassword)

	return nil
}

func (p *postStorageEu) Post(ctx context.Context, post string) (Post_id_obj, error) {
	logger := p.Logger(ctx)

	id := uuid.New()

	writeStartTimeMs := time.Now().UnixMilli()
	err := p.client.Set(ctx, id.String(), post, 0).Err()
	writePostDurationMs.Put(float64(time.Now().UnixMilli() - writeStartTimeMs))
	if err != nil {
		logger.Error("Error writing post!", "msg", err.Error())
		return Post_id_obj{}, err
	}

	logger.Debug("Post successfully stored!", "postId", id.String(), "post", post)

	return Post_id_obj{PostId: id.String(), WriteTime: writeStartTimeMs}, nil
}
