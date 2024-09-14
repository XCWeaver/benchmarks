package services

import (
	"context"
	"eu_deployment/pkg/model"
	"eu_deployment/pkg/storage"

	"github.com/TiagoMalhadas/xcweaver"

	"github.com/redis/go-redis/v9"
)

type HomeTimelineService interface {
	ReadHomeTimeline(ctx context.Context, reqID int64, userID string, start int64, stop int64) ([]model.Post, error)
}

type homeTimelineService struct {
	xcweaver.Implements[HomeTimelineService]
	xcweaver.WithConfig[homeTimelineServiceOptions]
	postStorageService xcweaver.Ref[PostStorageService]
	redisClient        *redis.Client
}

type homeTimelineServiceOptions struct {
	RedisAddr string `toml:"redis_address"`
	RedisPort int    `toml:"redis_port"`
	Region    string `toml:"region"`
}

func (h *homeTimelineService) Init(ctx context.Context) error {
	logger := h.Logger(ctx)
	h.redisClient = storage.RedisClient(h.Config().RedisAddr, h.Config().RedisPort)
	logger.Info("home timeline service running!", "region", h.Config().Region,
		"rabbitmq_addr", h.Config().RedisAddr, "rabbitmq_port", h.Config().RedisPort,
	)
	return nil
}

// readCachedTimeline is an helper function for reading timeline from redis with the same behavior as in the user timeline service
func (h *homeTimelineService) readCachedTimeline(ctx context.Context, userID string, start int64, stop int64) ([]string, error) {
	logger := h.Logger(ctx)

	result, err := h.redisClient.ZRevRange(ctx, userID, start, stop-1).Result()
	if err != nil {
		logger.Error("error reading home timeline from redis")
		return nil, err
	}

	var postIDs []string
	for _, result := range result {
		postIDs = append(postIDs, result)
	}
	return postIDs, nil
}

func (h *homeTimelineService) ReadHomeTimeline(ctx context.Context, reqID int64, userID string, start int64, stop int64) ([]model.Post, error) {
	logger := h.Logger(ctx)
	logger.Debug("entering ReadHomeTimeline", "req_id", reqID, "user_id", userID, "start", start, "stop", stop)
	if stop <= start || start < 0 {
		return []model.Post{}, nil
	}

	postIDs, err := h.readCachedTimeline(ctx, userID, start, stop)
	if err != nil {
		return nil, err
	}
	return h.postStorageService.Get().ReadPosts(ctx, reqID, postIDs)
}
