package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	sn_metrics "eu_deployment/pkg/metrics"
	"eu_deployment/pkg/model"
	"eu_deployment/pkg/storage"
	sn_trace "eu_deployment/pkg/trace"

	"github.com/ServiceWeaver/weaver"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ComposePostService interface {
	UploadCreator(ctx context.Context, reqID int64, creator model.Creator, new bool) error
	UploadText(ctx context.Context, reqID int64, text string, new bool) error
	UploadMedia(ctx context.Context, reqID int64, medias []model.Media, new bool) error
	UploadUniqueId(ctx context.Context, reqID int64, postID string, postType model.PostType, new bool) error
	UploadUrls(ctx context.Context, reqID int64, urls []model.URL, new bool) error
	UploadUserMentions(ctx context.Context, reqID int64, userMentions []model.UserMention, new bool) error
}

const NUM_COMPONENTS int = 6 // corresponds to the number of exposed methods
const REDIS_EXPIRE_TIME int = 12

type composePostService struct {
	weaver.Implements[ComposePostService]
	weaver.WithConfig[composePostServiceOptions]
	postStorageService  weaver.Ref[PostStorageService]
	userTimelineService weaver.Ref[UserTimelineService]
	redisClient         *redis.Client
	amqClientPool       *storage.RabbitMQClientPool
}

type composePostServiceOptions struct {
	RabbitMQAddr string   `toml:"rabbitmq_address"`
	RedisAddr    string   `toml:"redis_address"`
	RabbitMQPort int      `toml:"rabbitmq_port"`
	RedisPort    int      `toml:"redis_port"`
	Region       string   `toml:"region"`
	Regions      []string `toml:"regions"`
}

type MethodLabels struct {
	Caller    string // full calling component name
	Component string // full callee component name
	Method    string // callee component method's name
	Remote    bool   // Is this a remote call?
}

func (c *composePostService) Init(ctx context.Context) error {
	logger := c.Logger(ctx)
	var err error
	c.amqClientPool, err = storage.NewRabbitMQClientPool(ctx, c.Config().RabbitMQAddr, c.Config().RabbitMQPort, 0, 500)
	if err != nil {
		logger.Error("error initializing rabbitmq client pool", "msg", err.Error())
		return err
	}
	c.redisClient = storage.RedisClient(c.Config().RedisAddr, c.Config().RedisPort)
	logger.Info("compose post service running!", "region", c.Config().Region, "regions", c.Config().Regions,
		"rabbitmq_addr", c.Config().RabbitMQAddr, "rabbitmq_port", c.Config().RabbitMQPort,
		"redis_addr", c.Config().RedisAddr, "redis_port", c.Config().RedisPort,
	)
	return nil
}

func (c *composePostService) uploadComponent(ctx context.Context, reqID int64, fieldsValues ...interface{}) error {
	logger := c.Logger(ctx)
	reqIDStr := strconv.FormatInt(reqID, 10)
	cmds, err := c.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := pipe.HSet(ctx, reqIDStr, fieldsValues).Err()
		if err != nil {
			return err
		}
		err = pipe.HIncrBy(ctx, reqIDStr, "num_components", 1).Err()
		if err != nil {
			return err
		}
		err = pipe.Expire(ctx, reqIDStr, time.Second*time.Duration(REDIS_EXPIRE_TIME)).Err()
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil || len(cmds) != 3 /* sanity check */ {
		logger.Error("error writing component to redis", "fieldValues", fieldsValues, "msg", err.Error())
		return err
	}

	numComponents, err := cmds[1].(*redis.IntCmd).Result()
	if err != nil {
		logger.Error("error reading number of written components from redis", "fieldValues", fieldsValues, "msg", err.Error())
		return err
	}

	if numComponents == int64(NUM_COMPONENTS) {
		return c.composeAndUpload(ctx, reqID)
	}

	return nil
}

func (c *composePostService) uploadUpdatedComponent(ctx context.Context, reqID int64, fieldsValues ...interface{}) error {
	logger := c.Logger(ctx)
	reqIDStr := strconv.FormatInt(reqID, 10)
	cmds, err := c.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		err := pipe.HSet(ctx, reqIDStr, fieldsValues).Err()
		if err != nil {
			return err
		}
		err = pipe.HIncrBy(ctx, reqIDStr, "num_components", 1).Err()
		if err != nil {
			return err
		}
		err = pipe.Expire(ctx, reqIDStr, time.Second*time.Duration(REDIS_EXPIRE_TIME)).Err()
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil || len(cmds) != 3 /* sanity check  */ {
		logger.Error("error writing component to redis", "fieldValues", fieldsValues, "msg", err.Error())
		return err
	}

	numComponents, err := cmds[1].(*redis.IntCmd).Result()
	if err != nil {
		logger.Error("error reading number of written components from redis", "fieldValues", fieldsValues, "msg", err.Error())
		return err
	}

	if numComponents == int64(NUM_COMPONENTS) {
		return c.composeAndUpdate(ctx, reqID)
	}

	return nil
}

func (c *composePostService) UploadText(ctx context.Context, reqID int64, text string, new bool) error {
	logger := c.Logger(ctx)
	logger.Debug("entering UploadText", "text", text)
	textJSON, err := json.Marshal(text)
	if err != nil {
		logger.Error("error converting text to json", "text", text)
		return err
	}
	if !new {
		return c.uploadUpdatedComponent(ctx, reqID, "text", textJSON)
	}
	return c.uploadComponent(ctx, reqID, "text", textJSON)
}

func (c *composePostService) UploadMedia(ctx context.Context, reqID int64, medias []model.Media, new bool) error {
	logger := c.Logger(ctx)
	logger.Debug("entering UploadMedia", "medias", medias)
	mediasJSON, err := json.Marshal(medias)
	if err != nil {
		logger.Error("error converting medias to json", "medias", medias)
		return err
	}
	if !new {
		return c.uploadUpdatedComponent(ctx, reqID, "media", mediasJSON)
	}
	return c.uploadComponent(ctx, reqID, "media", mediasJSON)
}

func (c *composePostService) UploadUniqueId(ctx context.Context, reqID int64, postID string, postType model.PostType, new bool) error {
	logger := c.Logger(ctx)
	logger.Debug("entering UploadUniqueId", "post_id", postID, "post_type", postType)
	postIDJSON, err := json.Marshal(postID)
	if err != nil {
		logger.Error("error converting post id to json", "post_id", postID)
		return err
	}
	postTypeJSON, err := json.Marshal(postType)
	if err != nil {
		logger.Error("error converting medias to json", "post_type", postType)
		return err
	}
	if !new {
		return c.uploadUpdatedComponent(ctx, reqID, "post_id", postIDJSON, "post_type", postTypeJSON)
	}
	return c.uploadComponent(ctx, reqID, "post_id", postIDJSON, "post_type", postTypeJSON)
}

func (c *composePostService) UploadUrls(ctx context.Context, reqID int64, urls []model.URL, new bool) error {
	logger := c.Logger(ctx)
	logger.Debug("entering UploadUrls", "urls", urls)
	urlsJSON, err := json.Marshal(urls)
	if err != nil {
		logger.Error("error converting urls to json", "urls", urls)
		return err
	}
	if !new {
		return c.uploadUpdatedComponent(ctx, reqID, "urls", urlsJSON)
	}
	return c.uploadComponent(ctx, reqID, "urls", urlsJSON)
}

func (c *composePostService) UploadUserMentions(ctx context.Context, reqID int64, userMentions []model.UserMention, new bool) error {
	logger := c.Logger(ctx)
	logger.Debug("entering UploadUserMentions", "user_mentions", userMentions)
	userMentionsJSON, err := json.Marshal(userMentions)
	if err != nil {
		logger.Error("error converting user mentions to json", "user_mentions", userMentions)
		return err
	}
	if !new {
		return c.uploadUpdatedComponent(ctx, reqID, "user_mentions", userMentionsJSON)
	}
	return c.uploadComponent(ctx, reqID, "user_mentions", userMentionsJSON)
}

func (c *composePostService) UploadCreator(ctx context.Context, reqID int64, creator model.Creator, new bool) error {
	logger := c.Logger(ctx)
	logger.Debug("entering UploadCreator", "creator", creator)
	creatorJSON, err := json.Marshal(creator)
	if err != nil {
		logger.Error("error converting creator to json", "user_mentions", creatorJSON)
		return err
	}
	if !new {
		return c.uploadUpdatedComponent(ctx, reqID, "creator", creatorJSON)
	}
	return c.uploadComponent(ctx, reqID, "creator", creatorJSON)
}

func (c *composePostService) composeAndUpload(ctx context.Context, reqID int64) error {
	logger := c.Logger(ctx)
	logger.Debug("entering composeAndUpload", "reqid", reqID)

	// Recover from panic
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic:", r)
		}
	}()

	var text string
	var creator model.Creator
	var medias []model.Media
	var postID string
	var urls []model.URL
	var userMentions []model.UserMention
	var postType model.PostType

	var errs [7]error
	var wg sync.WaitGroup
	wg.Add(7)

	reqIDStr := strconv.FormatInt(reqID, 10)
	loadComponent := func(key string, value interface{}) error {
		logger.Debug("loading component", "reqid", reqIDStr, "key", key)
		cmd := c.redisClient.HGet(ctx, reqIDStr, key)
		if cmd == nil || cmd.Err() != nil {
			return cmd.Err()
		}
		result, err := cmd.Bytes()
		if err != nil {
			return err
		}
		return json.Unmarshal(result, &value)
	}

	go func() {
		defer wg.Done()
		errs[0] = loadComponent("text", &text)
	}()
	go func() {
		defer wg.Done()
		errs[1] = loadComponent("creator", &creator)
	}()
	go func() {
		defer wg.Done()
		errs[2] = loadComponent("media", &medias)
	}()
	go func() {
		defer wg.Done()
		errs[3] = loadComponent("post_id", &postID)
	}()
	go func() {
		defer wg.Done()
		errs[4] = loadComponent("urls", &urls)
	}()
	go func() {
		defer wg.Done()
		errs[5] = loadComponent("user_mentions", &userMentions)
	}()
	go func() {
		defer wg.Done()
		errs[6] = loadComponent("post_type", &postType)
	}()
	wg.Wait()
	logger.Debug("got all components from redis")

	for _, err := range errs {
		if err != nil {
			logger.Error("error reading from redis", "msg", err.Error())
			return err
		}
	}

	logger.Debug("parsing post data")
	timestamp := time.Now().UnixMilli()
	post := model.Post{
		PostID:       postID,
		ReqID:        reqID,
		Creator:      creator,
		Text:         text,
		UserMentions: userMentions,
		Media:        medias,
		URLs:         urls,
		Timestamp:    timestamp,
		PostType:     postType,
	}
	var userMentionIDs []string
	for _, mention := range userMentions {
		userMentionIDs = append(userMentionIDs, mention.UserID)
	}

	// --- Post Storage
	logger.Debug("remotely calling PostStorageService")

	regionLabel := sn_metrics.RegionLabel{Region: c.Config().Region}
	sn_metrics.ComposedPosts.Get(regionLabel).Inc()

	err := c.postStorageService.Get().StorePost(ctx, reqID, post)
	if err != nil {
		logger.Warn("error calling post storage service", "msg", err.Error())
		return err
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String("post_id", postID),
	)

	// --- Write Home Timeline
	logger.Debug("queueing message to rabbitmq")
	c.uploadHomeTimelineHelper(ctx, reqID, postID, creator.UserID, timestamp, userMentionIDs)

	// --- User Timeline
	logger.Debug("calling write user timeline")
	c.userTimelineService.Get().WriteUserTimeline(ctx, reqID, postID, post.Creator.UserID, timestamp)

	logger.Debug("done!")
	return nil
}

func (c *composePostService) composeAndUpdate(ctx context.Context, reqID int64) error {
	logger := c.Logger(ctx)
	logger.Debug("entering composeAndUpload", "reqid", reqID)

	var text string
	var creator model.Creator
	var medias []model.Media
	var postID string
	var urls []model.URL
	var userMentions []model.UserMention
	var postType model.PostType

	var errs [7]error
	var wg sync.WaitGroup
	wg.Add(7)

	reqIDStr := strconv.FormatInt(reqID, 10)
	loadComponent := func(key string, value interface{}) error {
		logger.Debug("loading component", "reqid", reqIDStr, "key", key)
		cmd := c.redisClient.HGet(ctx, reqIDStr, key)
		if cmd == nil || cmd.Err() != nil {
			return cmd.Err()
		}
		result, err := cmd.Bytes()
		if err != nil {
			return err
		}
		return json.Unmarshal(result, &value)
	}

	go func() {
		defer wg.Done()
		errs[0] = loadComponent("text", &text)
	}()
	go func() {
		defer wg.Done()
		errs[1] = loadComponent("creator", &creator)
	}()
	go func() {
		defer wg.Done()
		errs[2] = loadComponent("media", &medias)
	}()
	go func() {
		defer wg.Done()
		errs[3] = loadComponent("post_id", &postID)
	}()
	go func() {
		defer wg.Done()
		errs[4] = loadComponent("urls", &urls)
	}()
	go func() {
		defer wg.Done()
		errs[5] = loadComponent("user_mentions", &userMentions)
	}()
	go func() {
		defer wg.Done()
		errs[6] = loadComponent("post_type", &postType)
	}()
	wg.Wait()
	logger.Debug("got all components from redis")

	for _, err := range errs {
		if err != nil {
			logger.Error("error reading from redis", "msg", err.Error())
			return err
		}
	}

	logger.Debug("parsing post data")
	timestamp := time.Now().UnixMilli()
	post := model.Post{
		PostID:       postID,
		ReqID:        reqID,
		Creator:      creator,
		Text:         text,
		UserMentions: userMentions,
		Media:        medias,
		URLs:         urls,
		Timestamp:    timestamp,
		PostType:     postType,
	}
	var userMentionIDs []string
	for _, mention := range userMentions {
		userMentionIDs = append(userMentionIDs, mention.UserID)
	}

	// --- Post Storage
	logger.Debug("remotely calling PostStorageService")

	regionLabel := sn_metrics.RegionLabel{Region: c.Config().Region}
	sn_metrics.UpdatedPosts.Get(regionLabel).Inc()

	err := c.postStorageService.Get().UpdatePost(ctx, reqID, post)
	if err != nil {
		logger.Warn("error calling post storage service", "msg", err.Error())
		return err
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String("post_id", postID),
	)

	// --- Write Home Timeline
	logger.Debug("queueing message to rabbitmq")
	c.updateHomeTimelineHelper(ctx, reqID, postID, creator.UserID, timestamp, userMentionIDs)

	logger.Debug("done!")
	return nil
}

func (c *composePostService) uploadHomeTimelineHelper(ctx context.Context, reqID int64, postID string, userID string, timestamp int64, userMentionIDs []string) error {
	logger := c.Logger(ctx)

	ch, err := c.amqClientPool.Conn.Channel()
	defer c.amqClientPool.Push(ch)

	if err != nil {
		logger.Error("error getting rabbitmq client from pool", "msg", err.Error())
		panic(err)
	}
	err = ch.ExchangeDeclare("write-home-timeline", "topic", false, false, false, false, nil)
	if err != nil {
		logger.Error("error declaring exchange for rabbitmq", "msg", err.Error())
		// errors close the channel so we force the service to restart
		panic(err)
	}

	spanContext := trace.SpanContextFromContext(ctx)
	msg := model.Message{
		ReqID:          reqID,
		PostID:         postID,
		UserID:         userID,
		Timestamp:      timestamp,
		UserMentionIDs: userMentionIDs,
		// tracing
		SpanContext: sn_trace.BuildSpanContext(spanContext),
		// evaluation metrics
		NotificationSendTs: time.Now().UnixMilli(),
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		logger.Error("error converting rabbitmq message to json", "msg", err.Error())
		return err
	}

	span := trace.SpanFromContext(ctx)
	if trace.SpanContextFromContext(ctx).IsValid() {
		logger.Debug("valid span", "s", span.IsRecording())
	}

	amqMsg := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(msgJSON),
	}

	routingKey := "write-home-timeline-us"
	err = ch.PublishWithContext(ctx, "write-home-timeline", routingKey, false, false, amqMsg)
	if err != nil {
		logger.Error("error publishing to queue", "routing_key", routingKey, "err", err.Error())
	}
	regionLabel := sn_metrics.RegionLabel{Region: c.Config().Region}
	sn_metrics.SentNotifications.Get(regionLabel).Inc()

	//c.amqClientPool.Push(ch)

	span = trace.SpanFromContext(ctx)
	if trace.SpanContextFromContext(ctx).IsValid() {
		logger.Debug("valid span", "s", span.IsRecording())
	}
	return nil
}

func (c *composePostService) updateHomeTimelineHelper(ctx context.Context, reqID int64, postID string, userID string, timestamp int64, userMentionIDs []string) error {
	logger := c.Logger(ctx)

	ch, err := c.amqClientPool.Pop(ctx)
	defer c.amqClientPool.Push(ch)

	if err != nil {
		logger.Error("error getting rabbitmq client from pool", "msg", err.Error())
		panic(err)
	}
	err = ch.ExchangeDeclare("update-home-timeline", "topic", false, false, false, false, nil)
	if err != nil {
		logger.Error("error declaring exchange for rabbitmq", "msg", err.Error())
		// errors close the channel so we force the service to restart
		panic(err)
	}

	spanContext := trace.SpanContextFromContext(ctx)
	msg := model.Message{
		ReqID:          reqID,
		PostID:         postID,
		UserID:         userID,
		Timestamp:      timestamp,
		UserMentionIDs: userMentionIDs,
		// tracing
		SpanContext: sn_trace.BuildSpanContext(spanContext),
		// evaluation metrics
		NotificationSendTs: time.Now().UnixMilli(),
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		logger.Error("error converting rabbitmq message to json", "msg", err.Error())
		return err
	}

	span := trace.SpanFromContext(ctx)
	if trace.SpanContextFromContext(ctx).IsValid() {
		logger.Debug("valid span", "s", span.IsRecording())
	}

	amqMsg := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(msgJSON),
	}
	for _, region := range c.Config().Regions {
		routingKey := fmt.Sprintf("update-home-timeline-%s", region)
		err = ch.PublishWithContext(ctx, "update-home-timeline", routingKey, false, false, amqMsg)
		if err != nil {
			logger.Error("error publishing to queue", "routing_key", routingKey, "err", err.Error())
		}
		regionLabel := sn_metrics.RegionLabel{Region: c.Config().Region}
		sn_metrics.SentNotifications.Get(regionLabel).Inc()

	}
	c.amqClientPool.Push(ch)

	span = trace.SpanFromContext(ctx)
	if trace.SpanContextFromContext(ctx).IsValid() {
		logger.Debug("valid span", "s", span.IsRecording())
	}
	return nil
}
