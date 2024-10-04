package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	sn_metrics "eu_deployment/pkg/metrics"
	"eu_deployment/pkg/model"
	"eu_deployment/pkg/storage"

	"github.com/XCWeaver/xcweaver"
	"github.com/bradfitz/gomemcache/memcache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PostStorageService interface {
	StorePost(ctx context.Context, reqID int64, post model.Post) ([]byte, error)
	UpdatePost(ctx context.Context, reqID int64, post model.Post) ([]byte, error)
	ReadPost(ctx context.Context, reqID int64, postID string) (model.Post, error)
	ReadPosts(ctx context.Context, reqID int64, postIDs []string) ([]model.Post, error)
}

var _ xcweaver.NotRetriable = PostStorageService.StorePost

type postStorageServiceOptions struct {
	MongoDBAddr   string `toml:"mongodb_address"`
	MemCachedAddr string `toml:"memcached_address"`
	MongoDBPort   int    `toml:"mongodb_port"`
	MemCachedPort int    `toml:"memcached_port"`
	Region        string `toml:"region"`
}

type postStorageService struct {
	xcweaver.Implements[PostStorageService]
	xcweaver.WithConfig[postStorageServiceOptions]
	mongoClientPostStorage xcweaver.Antipode
	memCachedClient        *memcache.Client
	mongoClient            *mongo.Client
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
		"memcached_addr", p.Config().MemCachedAddr, "memcached_port", p.Config().MemCachedPort, "antipode", p.mongoClientPostStorage,
	)
	return nil
}

func (p *postStorageService) StorePost(ctx context.Context, reqID int64, post model.Post) ([]byte, error) {
	logger := p.Logger(ctx)
	logger.Info("entering StorePost", "reqid", reqID, "post", post)

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("poststorage_write_post_ts", time.Now().UnixMilli()),
	)

	postJSON, err := json.Marshal(post)
	if err != nil {
		logger.Debug("error converting post to JSON", "post", post, "msg", err.Error())
		return nil, err
	}

	writePostStartMs := time.Now().UnixMilli()
	ctx, err = p.mongoClientPostStorage.Write(ctx, "posts", post.PostID, string(postJSON))

	/*collection := p.mongoClient.Database("post-storage").Collection("posts")
	r, err := collection.InsertOne(ctx, post)*/
	if err != nil {
		logger.Error("error writing post", "msg", err.Error())
		return nil, err
	}
	logger.Debug("write post done!", "key", post.PostID, "post", string(postJSON))
	regionLabel := sn_metrics.RegionLabel{Region: p.Config().Region}
	logger.Debug("before write post metric 1", "region_label", regionLabel)
	sn_metrics.WritePostDurationMs.Get(regionLabel)
	logger.Debug("before write post metric 2", "region_label", regionLabel)
	sn_metrics.WritePostDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - writePostStartMs))
	//logger.Debug("inserted post", "objectid", r.InsertedID)

	lineage, err := xcweaver.GetLineage(ctx)
	if err != nil {
		logger.Error("error getting lineage from context", "msg", err.Error())
		return nil, err
	}

	return lineage, nil
}

func (p *postStorageService) UpdatePost(ctx context.Context, reqID int64, post model.Post) ([]byte, error) {
	logger := p.Logger(ctx)
	logger.Info("entering UpdatePost", "reqid", reqID, "post", post)

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Int64("poststorage_update_post_ts", time.Now().UnixMilli()),
	)

	postJSON, err := json.Marshal(post)
	if err != nil {
		logger.Debug("error converting post to JSON", "post", post, "msg", err.Error())
		return nil, err
	}
	postStr := string(postJSON)

	updatePostStartMs := time.Now().UnixMilli()
	ctx, err = p.mongoClientPostStorage.Write(ctx, "posts", post.PostID, postStr)

	if err != nil {
		logger.Error("error writing post", "msg", err.Error())
		return nil, err
	}
	logger.Debug("update post done!", "key", post.PostID, "post", postStr)
	regionLabel := sn_metrics.RegionLabel{Region: p.Config().Region}
	logger.Debug("before update post metric 1", "region_label", regionLabel)
	sn_metrics.UpdatePostDurationMs.Get(regionLabel)
	logger.Debug("before update post metric 2", "region_label", regionLabel)
	sn_metrics.UpdatePostDurationMs.Get(regionLabel).Put(float64(time.Now().UnixMilli() - updatePostStartMs))

	lineage, err := xcweaver.GetLineage(ctx)
	if err != nil {
		logger.Error("error getting lineage from context", "msg", err.Error())
		return nil, err
	}
	return lineage, nil
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
		result, lineage, err := p.mongoClientPostStorage.Read(ctx, "posts", postID)
		if err != nil {
			logger.Error("error reading post from mongo", "msg", err.Error())
			return post, err
		}
		logger.Debug("postStorage service | lineage and message successfully read!", "region", p.Config().Region, "lineage", lineage, "message", result)
		err = json.Unmarshal([]byte(result), &post)
		if err != nil {
			errMsg := fmt.Sprintf("post_id: %s not found in mongodb", postID)
			logger.Warn(errMsg)
			return post, fmt.Errorf(errMsg)
		}
	}

	return post, nil
}

// To-Do
// Use Antipode
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
