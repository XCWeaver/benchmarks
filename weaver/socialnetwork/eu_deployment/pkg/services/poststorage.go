package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	sn_metrics "eu_deployment/pkg/metrics"
	"eu_deployment/pkg/model"
	"eu_deployment/pkg/storage"

	"github.com/ServiceWeaver/weaver"
	"github.com/bradfitz/gomemcache/memcache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type contextKey string

type PostStorageService interface {
	StorePost(ctx context.Context, reqID int64, post model.Post) error
	UpdatePost(ctx context.Context, reqID int64, post model.Post) error
	ReadPost(ctx context.Context, reqID int64, postID string) (model.Post, error)
	ReadPosts(ctx context.Context, reqID int64, postIDs []string) ([]model.Post, error)
}

var _ weaver.NotRetriable = PostStorageService.StorePost

type postStorageServiceOptions struct {
	MongoDBAddr   string `toml:"mongodb_address"`
	MemCachedAddr string `toml:"memcached_address"`
	MongoDBPort   int    `toml:"mongodb_port"`
	MemCachedPort int    `toml:"memcached_port"`
	Region        string `toml:"region"`
}

type postStorageService struct {
	weaver.Implements[PostStorageService]
	weaver.WithConfig[postStorageServiceOptions]
	mongoClient     *mongo.Client
	memCachedClient *memcache.Client
}

func (p *postStorageService) Init(ctx context.Context) error {
	logger := p.Logger(ctx)
	var err error
	p.mongoClient, err = storage.MongoDBClient(ctx, p.Config().MongoDBAddr, p.Config().MongoDBPort)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	p.memCachedClient = storage.MemCachedClient(p.Config().MemCachedAddr, p.Config().MemCachedPort)
	if p.memCachedClient == nil {
		errMsg := "error connecting to memcached"
		logger.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	logger.Info("post storage service running!", "region", p.Config().Region,
		"mongodb_addr", p.Config().MongoDBAddr, "mongodb_port", p.Config().MongoDBPort,
		"memcached_addr", p.Config().MemCachedAddr, "memcached_port", p.Config().MemCachedPort,
	)
	return nil
}

func (p *postStorageService) StorePost(ctx context.Context, reqID int64, post model.Post) error {
	logger := p.Logger(ctx)
	logger.Info("entering StorePost", "reqid", reqID, "post", post)

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("poststorage_write_post_ts", time.Now().UnixMilli()),
	)
	writePostStartMs := time.Now().UnixMilli()

	collection := p.mongoClient.Database("post-storage").Collection("posts")
	r, err := collection.InsertOne(ctx, post)
	if err != nil {
		logger.Error("error writing post", "msg", err.Error())
	}
	regionLabel := sn_metrics.RegionLabel{Region: p.Config().Region}
	logger.Debug("before write post metric 1", "region_label", regionLabel)
	sn_metrics.WritePostDurationMs.Get(regionLabel)
	logger.Debug("before write post metric 2", "region_label", regionLabel)
	sn_metrics.WritePostDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - writePostStartMs))
	logger.Debug("inserted post", "objectid", r.InsertedID)

	return nil
}

func (p *postStorageService) UpdatePost(ctx context.Context, reqID int64, post model.Post) error {
	logger := p.Logger(ctx)
	logger.Info("entering UpdatePost", "reqid", reqID, "post", post)

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("poststorage_write_post_ts", time.Now().UnixMilli()),
	)
	updatePostStartMs := time.Now().UnixMilli()

	collection := p.mongoClient.Database("post-storage").Collection("posts")
	filter := bson.D{{"post_id", post.PostID}}
	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return errors.New("Post not found!")
	} else if result.Err() != nil {
		return result.Err()
	}

	update := bson.D{{"$set", bson.D{{"req_id", post.ReqID}, {"creator", post.Creator}, {"text", post.Text}, {"user_mentions", post.UserMentions},
		{"media", post.Media}, {"urls", post.URLs}, {"timestamp", post.Timestamp}, {"posttype", post.PostType}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	//delete post from cache
	err = p.memCachedClient.Delete(post.PostID)
	if err != nil && err != memcache.ErrCacheMiss {
		// error reading cache
		logger.Error("error reading post from cache", "msg", err.Error())
		return err
	}

	regionLabel := sn_metrics.RegionLabel{Region: p.Config().Region}
	logger.Debug("before update post metric 1", "region_label", regionLabel)
	sn_metrics.UpdatePostDurationMs.Get(regionLabel)
	logger.Debug("before update post metric 2", "region_label", regionLabel)
	sn_metrics.UpdatePostDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - updatePostStartMs))
	logger.Debug("post updated")

	return nil
}

func (p *postStorageService) ReadPost(ctx context.Context, reqID int64, postID string) (model.Post, error) {
	logger := p.Logger(ctx)
	logger.Info("entering ReadPost", "req_id", reqID, "post_id", postID)

	var post model.Post
	item, err := p.memCachedClient.Get(postID)

	if err != nil && err != memcache.ErrCacheMiss {
		// error reading cache
		logger.Error("error reading post from cache", "msg", err.Error())
		return post, err
	}
	if err == nil {
		// post found in cache
		err := json.Unmarshal(item.Value, &post)
		if err != nil {
			logger.Error("error parsing post from cache result", "msg", err.Error())
			return post, err
		}
	} else {
		// post does not exist in cache
		// so we get it from db
		collection := p.mongoClient.Database("post-storage").Collection("posts")
		filter := bson.D{
			{Key: "post_id", Value: postID},
		}
		result := collection.FindOne(ctx, filter)
		if result.Err() != nil {
			return post, err
		}
		err = result.Decode(&post)
		if err != nil {
			errMsg := fmt.Sprintf("post_id: %s not found in mongodb", postID)
			logger.Warn(errMsg)
			return post, fmt.Errorf(errMsg)
		}
	}
	return post, nil
}

func (p *postStorageService) ReadPosts(ctx context.Context, reqID int64, postIDs []string) ([]model.Post, error) {
	logger := p.Logger(ctx)
	logger.Info("entering ReadPosts", "req_id", reqID, "post_ids", postIDs)

	if len(postIDs) == 0 {
		return []model.Post{}, nil
	}

	postIDsNotCached := make(map[string]bool)
	for _, pid := range postIDs {
		postIDsNotCached[pid] = true
	}

	var keys []string
	for _, pid := range postIDs {
		keys = append(keys, pid)
	}
	result, err := p.memCachedClient.GetMulti(keys)
	if err != nil {
		logger.Error("error reading keys from memcached", "msg", err.Error())
		return nil, err
	}
	posts := []model.Post{}
	for _, key := range keys {
		if val, ok := result[key]; ok {
			var cachedPost model.Post
			err := json.Unmarshal(val.Value, &cachedPost)
			if err != nil {
				logger.Error("error parsing ids from memcached result", "msg", err.Error())
				return nil, err
			}
			posts = append(posts, cachedPost)
		}
	}

	for _, post := range posts {
		delete(postIDsNotCached, post.PostID)
	}
	if len(postIDsNotCached) != 0 {
		collection := p.mongoClient.Database("post-storage").Collection("posts")

		queryPostIDArray := bson.A{}
		for id := range postIDsNotCached {
			queryPostIDArray = append(queryPostIDArray, id)
		}
		filter := bson.D{
			{Key: "post_id", Value: bson.D{
				{Key: "$in", Value: queryPostIDArray},
			}},
		}
		cur, err := collection.Find(ctx, filter)
		if err != nil {
			logger.Error("error reading posts from mongodb", "msg", err.Error())
			return nil, err
		}

		exists := cur.TryNext(ctx)
		if exists {
			var postsToCache []model.Post
			err = cur.All(ctx, &postsToCache)
			if err != nil {
				logger.Error("error parsing new posts from mongodb", "msg", err.Error())
				return nil, err
			}
			posts = append(posts, postsToCache...)

			var wg sync.WaitGroup
			for _, newPost := range postsToCache {
				wg.Add(1)

				go func(newPost model.Post) error {
					defer wg.Done()
					postJson, err := json.Marshal(newPost)
					if err != nil {
						logger.Error("error converting post to json", "post", newPost)
						return err
					}
					p.memCachedClient.Set(&memcache.Item{Key: newPost.PostID, Value: postJson})
					return nil
				}(newPost)
			}
			wg.Wait()
		}
	}
	return posts, nil
}
