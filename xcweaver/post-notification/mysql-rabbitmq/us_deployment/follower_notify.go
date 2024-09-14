package main

import (
	"context"

	"github.com/TiagoMalhadas/xcweaver"
)

// Follower_Notify component.
type Follower_Notify interface {
	Follower_Notify(context.Context, Post_id_obj, int) error
}

// Implementation of the Follower_Notify component.
type follower_Notify struct {
	xcweaver.Implements[Follower_Notify]
	post_storage xcweaver.Ref[Post_storage_america]
	clientRedis  xcweaver.Antipode
}

func (f *follower_Notify) Init(ctx context.Context) error {
	logger := f.Logger(ctx)

	logger.Info("follower notify service running!")
	return nil
}

func (f *follower_Notify) Follower_Notify(ctx context.Context, postId Post_id_obj, userId int) error {
	logger := f.Logger(ctx)

	logger.Debug("calling barrier!")
	err := f.clientRedis.Barrier(ctx)
	if err != nil {
		logger.Error("error on barrier", "msg", err.Error())
		return err
	}
	logger.Debug("barrier executed successfully!")

	post, err := f.post_storage.Get().GetPost(ctx, postId)
	if err != nil {
		return err
	}

	logger.Debug("Post read successfully!", "post", post)

	return nil
}
