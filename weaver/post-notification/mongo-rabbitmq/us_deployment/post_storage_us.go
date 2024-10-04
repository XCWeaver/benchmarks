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

// PostStorageUs component.
type PostStorageUs interface {
	GetPost(context.Context, Post_id_obj) (string, error)
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
	GetInconsistencies(ctx context.Context) (int, error)
	Reset(ctx context.Context) error
}

type Post_id_obj struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}

// Implementation of the PostStorageUs component.
type postStorageUs struct {
	weaver.Implements[PostStorageUs]
	weaver.WithConfig[postStorageUsOptions]
	client                  *mongo.Client
	mu                      sync.Mutex
	consistencyWindowValues []float64
	muInc                   sync.Mutex
	inconsistencies         int
}

type postStorageUsOptions struct {
	MongoAddr string `toml:"mongo_address"`
	MongoPort string `toml:"mongo_port"`
}

type Post struct {
	Key  string `bson:"key"`
	Post string `bson:"post"`
}

func (p *postStorageUs) Init(ctx context.Context) error {
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

func (p *postStorageUs) GetPost(ctx context.Context, postId Post_id_obj) (string, error) {
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
		p.muInc.Lock()
		p.inconsistencies += 1
		p.muInc.Unlock()
		logger.Error("post not found")
		return "post not found", err
	} else if err != nil {
		return "", err
	}

	return result.Post, nil
}

func (p *postStorageUs) GetConsistencyWindowValues(ctx context.Context) ([]float64, error) {
	logger := p.Logger(ctx)
	logger.Debug("entering GetConsistencyWindowValues")
	p.mu.Lock()
	values := p.consistencyWindowValues
	p.mu.Unlock()
	return values, nil
}

func (p *postStorageUs) GetInconsistencies(ctx context.Context) (int, error) {
	logger := p.Logger(ctx)
	logger.Debug("entering GetInconsistencies")
	p.muInc.Lock()
	inconsistencies := p.inconsistencies
	p.muInc.Unlock()
	return inconsistencies, nil
}

func (p *postStorageUs) Reset(ctx context.Context) error {
	logger := p.Logger(ctx)
	logger.Debug("entering Reset")
	p.inconsistencies = 0
	p.consistencyWindowValues = []float64{}
	return nil
}
