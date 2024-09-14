package main

import (
	"context"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Post_storage component.
type Post_storage_europe interface {
	Post(context.Context, string) (Post_id_obj, error)
}

type Post_id_obj struct {
	weaver.AutoMarshal
	PostId    string
	WriteTime int64
}

// Implementation of the Post_storage component.
type post_storage_europe struct {
	weaver.Implements[Post_storage_europe]
	weaver.WithConfig[post_storage_europeOptions]
	client *mongo.Client
}

type post_storage_europeOptions struct {
	MongoAddr string `toml:"mongo_address"`
	MongoPort string `toml:"mongo_port"`
}

type Post struct {
	Key  string `bson:"key"`
	Post string `bson:"post"`
}

func (p *post_storage_europe) Init(ctx context.Context) error {
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

func (p *post_storage_europe) Post(ctx context.Context, post string) (Post_id_obj, error) {
	logger := p.Logger(ctx)

	id := uuid.New()

	writeStartTimeMs := time.Now().UnixMilli()

	collection := p.client.Database("post-storage").Collection("posts")

	postObj := Post{
		Key:  id.String(),
		Post: post,
	}

	_, err := collection.InsertOne(ctx, postObj)
	writePostDurationMs.Put(float64(time.Now().UnixMilli() - writeStartTimeMs))
	if err != nil {
		logger.Error("Error writing post!", "msg", err.Error())
		return Post_id_obj{}, err
	}

	logger.Debug("Post successfully stored!", "postId", id.String(), "post", post)

	return Post_id_obj{PostId: id.String(), WriteTime: writeStartTimeMs}, nil
}
