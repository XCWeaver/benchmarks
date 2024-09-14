package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	sn_metrics "us_deployment/pkg/metrics"
	"us_deployment/pkg/model"
	"us_deployment/pkg/storage"

	"github.com/TiagoMalhadas/xcweaver"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type WriteHomeTimelineService interface {
	//WriteHomeTimeline(ctx context.Context, msg model.Message) error
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
	ReadNotifications(ctx context.Context) error
}

type writeHomeTimelineServiceOptions struct {
	RedisAddr  string `toml:"redis_address"`
	RedisPort  int    `toml:"redis_port"`
	NumWorkers int    `toml:"num_workers"`
	Region     string `toml:"region"`
}

type writeHomeTimelineService struct {
	xcweaver.Implements[WriteHomeTimelineService]
	xcweaver.WithConfig[writeHomeTimelineServiceOptions]
	socialGraphService      xcweaver.Ref[SocialGraphService]
	redisClient             *redis.Client
	rabbitClientWriteHomeTL xcweaver.Antipode
	mongoClientWriteHomeTL  xcweaver.Antipode
	mu                      sync.Mutex
	consistencyWindowValues []float64
}

func (w *writeHomeTimelineService) Init(ctx context.Context) error {
	logger := w.Logger(ctx)

	w.redisClient = storage.RedisClient(w.Config().RedisAddr, w.Config().RedisPort)

	/*var wg sync.WaitGroup
	wg.Add(w.Config().NumWorkers)
	for i := 1; i <= w.Config().NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			err := w.workerThread(ctx, i)
			logger.Error("error in worker thread", "msg", err.Error())
		}(i)
	}*/

	logger.Info("write home timeline service running!", "region", w.Config().Region, "n_workers", w.Config().NumWorkers,
		"redis_addr", w.Config().RedisAddr, "redis_port", w.Config().RedisPort,
	)

	//wg.Wait()
	return nil
}

func (w *writeHomeTimelineService) WriteHomeTimeline(ctx context.Context, msg model.Message) error {
	logger := w.Logger(ctx)

	span := trace.SpanFromContext(ctx)
	if trace.SpanContextFromContext(ctx).IsValid() {
		logger.Debug("valid span", "s", span.IsRecording(), "ctx", ctx.Value("TEST"))
	}

	beginningRead := time.Now().UnixMilli()
	result, _, err := w.mongoClientWriteHomeTL.Read(ctx, "posts", msg.PostID)
	regionLabel := sn_metrics.RegionLabel{Region: w.Config().Region}
	sn_metrics.ReadPostDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - beginningRead))
	consistencyWindow := float64(time.Now().UnixMilli() - msg.Timestamp)
	sn_metrics.ConsistencyWindow.Get(regionLabel).Put(consistencyWindow)
	w.mu.Lock()
	w.consistencyWindowValues = append(w.consistencyWindowValues, consistencyWindow)
	w.mu.Unlock()
	if err != nil {
		if err == xcweaver.ErrNotFound {
			trace.SpanFromContext(ctx).SetAttributes(
				attribute.Bool("poststorage_consistent_read", false),
			)
			logger.Debug("inconsistency!")
			sn_metrics.Inconsistencies.Get(regionLabel).Inc()
			return nil
		} else {
			logger.Error("error reading post from mongodb", "msg", err.Error())
			return err
		}
	}

	var post model.Post
	err = json.Unmarshal([]byte(result), &post)
	if err != nil {
		errMsg := fmt.Sprintf("post_id: %s not found in mongodb", msg.PostID)
		logger.Warn(errMsg)
		return fmt.Errorf(errMsg)
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Bool("poststorage_consistent_read", true),
	)

	logger.Debug("found post! :)", "post_id", post.PostID, "text", post.Text)

	followersID, err := w.socialGraphService.Get().GetFollowers(ctx, msg.ReqID, msg.UserID)
	if err != nil {
		logger.Error("error getting followers from social graph service")
		return err
	}

	logger.Debug("got followers to write to their hometimeline", "num", len(followersID))
	uniqueIDs := make(map[string]bool, 0)
	for _, followerID := range followersID {
		uniqueIDs[followerID] = true
	}
	for _, userMentionID := range msg.UserMentionIDs {
		uniqueIDs[userMentionID] = true
	}
	value := redis.Z{
		Member: msg.PostID,
		Score:  float64(msg.Timestamp),
	}
	_, err = w.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for id := range uniqueIDs {
			err = w.redisClient.ZAddNX(ctx, id, value).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})
	logger.Debug("leaving write home timeline")
	return nil
}

// onReceivedWorker adds the post to all the post's subscribed users (followers, mentioned users, etc)
func (w *writeHomeTimelineService) onReceivedWorker(ctx context.Context, body []byte) error {
	logger := w.Logger(ctx)

	var msg model.Message
	err := json.Unmarshal(body, &msg)
	if err != nil {
		logger.Error("error parsing json message", "msg", err.Error())
		return err
	}
	regionLabel := sn_metrics.RegionLabel{Region: w.Config().Region}
	sn_metrics.ReceivedNotifications.Get(regionLabel).Add(1)
	sn_metrics.QueueDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - msg.NotificationSendTs))
	logger.Debug("received rabbitmq message", "post_id", msg.PostID, "msg", msg)

	err = w.mongoClientWriteHomeTL.Barrier(ctx)
	if err != nil {
		logger.Error("error on barrier", "msg", err.Error())
		return err
	}

	/*spanContext, err := sn_trace.ParseSpanContext(msg.SpanContext)
	if err != nil {
		logger.Error("error parsing span context", "workerid", workerid, "msg", err.Error())
		return err
	}

	ctx = trace.ContextWithRemoteSpanContext(ctx, spanContext)*/

	return w.WriteHomeTimeline(ctx, msg)
}

func (w *writeHomeTimelineService) ReadNotifications(ctx context.Context) error {
	logger := w.Logger(ctx)

	var forever chan struct{}
	var stop chan struct{}

	go func() {
		defer close(stop)

		msgs, err := w.rabbitClientWriteHomeTL.Consume(ctx, "write-home-timeline", "write-home-timeline-us", stop)
		if err != nil {
			return
		}

		for d := range msgs {
			go func(msg xcweaver.AntipodeObject) {
				ctx, err = xcweaver.Transfer(ctx, msg.Lineage)
				if err != nil {
					logger.Error("error transfering the lineage to context", "msg", err.Error())
					return
				}

				err = w.onReceivedWorker(ctx, []byte(msg.Version))
				if err != nil {
					logger.Warn("error in worker thread", "msg", err.Error())
				}
			}(d)
		}
	}()

	<-forever
	return nil
}

func (w *writeHomeTimelineService) GetConsistencyWindowValues(ctx context.Context) ([]float64, error) {
	w.mu.Lock()
	values := w.consistencyWindowValues
	w.mu.Unlock()
	return values, nil
}
