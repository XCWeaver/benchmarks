package main

import (
	"context"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Post_storage_america component.
type Post_storage_america interface {
	GetPost(context.Context, Post_id_obj) (string, error)
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
}

type Post_id_obj struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}

// Implementation of the Post_storage_america component.
type post_storage_america struct {
	weaver.Implements[Post_storage_america]
	weaver.WithConfig[post_storage_usOptions]
	client                  *mongo.Client
	mu                      sync.Mutex
	consistencyWindowValues []float64
}

type post_storage_usOptions struct {
	MongoAddr string `toml:"mongo_address"`
	MongoPort string `toml:"mongo_port"`
}

type Post struct {
	Key  string `bson:"key"`
	Post string `bson:"post"`
}

func (p *post_storage_america) Init(ctx context.Context) error {
	logger := p.Logger(ctx)

	var err error
	clientOptions := options.Client().ApplyURI("mongodb://" + p.Config().MongoAddr + ":" + p.Config().MongoPort + "/?directConnection=true")
	p.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("error conecting to mongoDB", "msg", err.Error())
		return err
	}

	logger.Info("post storage service at eu running!", "mongo host", p.Config().MongoAddr, "port", p.Config().MongoPort)

	return nil
}

func (p *post_storage_america) GetPost(ctx context.Context, postId Post_id_obj) (string, error) {
	logger := p.Logger(ctx)

	logger.Debug("Reading post!", "postId", postId.PostId)
	startTimeMs := time.Now().UnixMilli()

	filter := bson.D{{"key", postId.PostId}}

	var result Post
	err := p.client.Database("post-storage").Collection("posts").FindOne(ctx, filter).Decode(&result)
	readPostDurationMs.Put(float64(time.Now().UnixMilli() - startTimeMs))
	consistencyWindowMs := float64(time.Now().UnixMilli() - postId.WriteTime)
	consistencyWindow.Put(consistencyWindowMs)
	p.mu.Lock()
	p.consistencyWindowValues = append(p.consistencyWindowValues, consistencyWindowMs)
	p.mu.Unlock()
	if err == mongo.ErrNoDocuments {
		inconsistencies.Inc()
		logger.Error("post not found")
		return "post not found", err
	} else if err != nil {
		return "", err
	}

	return result.Post, nil
}

func (p *post_storage_america) GetConsistencyWindowValues(ctx context.Context) ([]float64, error) {
	logger := p.Logger(ctx)
	logger.Debug("entering GetConsistencyWindowValues")
	p.mu.Lock()
	values := p.consistencyWindowValues
	p.mu.Unlock()
	return values, nil
}
