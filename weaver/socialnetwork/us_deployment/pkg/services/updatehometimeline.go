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
	sn_trace "us_deployment/pkg/trace"

	"github.com/ServiceWeaver/weaver"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UpdateHomeTimelineService interface {
	UpdateHomeTimeline(ctx context.Context, msg model.Message) error
}

type updateHomeTimelineServiceOptions struct {
	RabbitMQAddr string `toml:"rabbitmq_address"`
	MongoDBAddr  string `toml:"mongodb_address"`
	RedisAddr    string `toml:"redis_address"`
	RabbitMQPort int    `toml:"rabbitmq_port"`
	MongoDBPort  int    `toml:"mongodb_port"`
	RedisPort    int    `toml:"redis_port"`
	NumWorkers   int    `toml:"num_workers"`
	Region       string `toml:"region"`
}

type updateHomeTimelineService struct {
	weaver.Implements[UpdateHomeTimelineService]
	weaver.WithConfig[updateHomeTimelineServiceOptions]
	socialGraphService weaver.Ref[SocialGraphService]
	mongoClient        *mongo.Client
	redisClient        *redis.Client
	amqClientPool      *storage.RabbitMQClientPool
}

func (w *updateHomeTimelineService) Init(ctx context.Context) error {
	logger := w.Logger(ctx)
	var err error
	w.mongoClient, err = storage.MongoDBClient(ctx, w.Config().MongoDBAddr, w.Config().MongoDBPort)
	if err != nil {
		logger.Error("error initializing mongodb client", "msg", err.Error())
		return err
	}
	w.redisClient = storage.RedisClient(w.Config().RedisAddr, w.Config().RedisPort)
	w.amqClientPool, err = storage.NewRabbitMQClientPool(ctx, w.Config().RabbitMQAddr, w.Config().RabbitMQPort, 0, 500)
	if err != nil {
		logger.Error("error initializing rabbitmq client pool", "msg", err.Error())
		return err
	}

	var wg sync.WaitGroup
	wg.Add(w.Config().NumWorkers)
	for i := 1; i <= w.Config().NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			err := w.workerThread(ctx, i)
			logger.Error("error in worker thread", "msg", err.Error())
		}(i)
	}

	logger.Info("update home timeline service running!", "region", w.Config().Region, "n_workers", w.Config().NumWorkers,
		"rabbitmq_addr", w.Config().RabbitMQAddr, "rabbitmq_port", w.Config().RabbitMQPort,
		"mongodb_addr", w.Config().MongoDBAddr, "mongodb_port", w.Config().MongoDBPort,
		"redis_addr", w.Config().RedisAddr, "redis_port", w.Config().RedisPort,
	)
	wg.Wait()
	return nil
}

func (w *updateHomeTimelineService) UpdateHomeTimeline(ctx context.Context, msg model.Message) error {
	logger := w.Logger(ctx)

	span := trace.SpanFromContext(ctx)
	if trace.SpanContextFromContext(ctx).IsValid() {
		logger.Debug("valid span", "s", span.IsRecording(), "ctx", ctx.Value("TEST"))
	}

	regionLabel := sn_metrics.RegionLabel{Region: w.Config().Region}
	sn_metrics.QueueDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - msg.NotificationSendTs))

	beginningRead := time.Now().UnixMilli()
	db := w.mongoClient.Database("post-storage")
	collection := db.Collection("posts")

	var post model.Post
	filter := bson.D{{Key: "post_id", Value: msg.PostID}}
	err := collection.FindOne(ctx, filter, nil).Decode(&post)
	sn_metrics.ReadPostDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - beginningRead))
	if err != nil {
		if err == mongo.ErrNoDocuments {
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

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Bool("poststorage_consistent_read", true),
	)

	//Used to test Antipode
	if post.Text != "updatedText" {
		trace.SpanFromContext(ctx).SetAttributes(
			attribute.Bool("poststorage_consistent_read", false),
		)
		logger.Debug("inconsistency!")
		sn_metrics.UpdateInconsistencies.Get(regionLabel).Inc()
		return nil
	}

	logger.Debug("found post! :)", "post_id", post.PostID, "text", post.Text)

	/*followersID, err := w.socialGraphService.Get().GetFollowers(ctx, msg.ReqID, msg.UserID)
	if err != nil {
		logger.Error("error getting followers from social graph service")
		return err
	}

	logger.Debug("got followers to write to their hometimeline", "num", len(followersID))
	uniqueIDs := make(map[int64]bool, 0)
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
			idStr := strconv.FormatInt(id, 10)
			err = w.redisClient.ZAddNX(ctx, idStr, value).Err()
			if err != nil {
				return err
			}
		}
		return nil
	})*/
	logger.Debug("leaving write home timeline")
	return nil
}

// onReceivedWorker adds the post to all the post's subscribed users (followers, mentioned users, etc)
func (w *updateHomeTimelineService) onReceivedWorker(ctx context.Context, workerid int, body []byte) error {
	logger := w.Logger(ctx)

	var msg model.Message
	err := json.Unmarshal(body, &msg)
	if err != nil {
		logger.Error("error parsing json message", "workerid", workerid, "msg", err.Error())
		return err
	}
	regionLabel := sn_metrics.RegionLabel{Region: w.Config().Region}
	sn_metrics.ReceivedNotifications.Get(regionLabel).Add(1)
	logger.Debug("received rabbitmq message", "workerid", workerid, "post_id", msg.PostID, "msg", msg)

	spanContext, err := sn_trace.ParseSpanContext(msg.SpanContext)
	if err != nil {
		logger.Error("error parsing span context", "workerid", workerid, "msg", err.Error())
		return err
	}

	ctx = trace.ContextWithRemoteSpanContext(ctx, spanContext)

	/* span := trace.SpanFromContext(ctx)
	ctx, span = s.tracer.Start(ctx, "services.WriteHomeTimelineService.WriteHomeTimeline", trace.WithSpanKind(trace.SpanKindInternal))
	defer func() {
		span.End()
	}() */

	return w.UpdateHomeTimeline(ctx, msg)
}

func (w *updateHomeTimelineService) workerThread(ctx context.Context, workerid int) error {
	logger := w.Logger(ctx)

	ch, err := w.amqClientPool.Pop(ctx)
	defer w.amqClientPool.Push(ch)

	if err != nil {
		logger.Error("error getting rabbitmq client from pool", "msg", err.Error())
		panic(err)
	}

	err = ch.ExchangeDeclare("update-home-timeline", "topic", false, false, false, false, nil)
	if err != nil {
		logger.Error("error declaring exchange for rabbitmq", "workerid", workerid, "msg", err.Error())
		return err
	}
	routingKey := fmt.Sprintf("update-home-timeline-%s", w.Config().Region)
	_, err = ch.QueueDeclare(routingKey, false, false, false, false, nil)
	if err != nil {
		logger.Error("error declaring queue for rabbitmq", "workerid", workerid, "msg", err.Error())
		return err
	}

	err = ch.QueueBind(routingKey, routingKey, "update-home-timeline", false, nil)
	if err != nil {
		logger.Error("error binding queue for rabbitmq", "workerid", workerid, "msg", err.Error())
		return err
	}

	msgs, err := ch.Consume(routingKey, "", true, false, false, false, nil)
	if err != nil {
		logger.Error("error consuming queue", "workerid", workerid, "msg", err.Error())
		return err
	}

	var forever chan struct{}
	go func() {
		for msg := range msgs {
			err = w.onReceivedWorker(ctx, workerid, msg.Body)
			if err != nil {
				logger.Warn("error in worker thread", "msg", err.Error())
			}
		}
	}()
	<-forever
	return nil
}
