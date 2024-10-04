package main

import (
	"context"

	"github.com/XCWeaver/xcweaver"
)

// FollowerNotify component.
type FollowerNotify interface {
	Follower_Notify(context.Context, Post_id_obj, int) error
}

// Implementation of the FollowerNotify component.
type followerNotify struct {
	xcweaver.Implements[FollowerNotify]
	post_storage xcweaver.Ref[PostStorageUS]
	client       xcweaver.Antipode
}

func (f *followerNotify) Init(ctx context.Context) error {
	logger := f.Logger(ctx)

	logger.Info("follower notify service running!")
	return nil
}

func (f *followerNotify) Follower_Notify(ctx context.Context, postId Post_id_obj, userId int) error {
	logger := f.Logger(ctx)

	logger.Debug("calling barrier!")
	err := f.client.Barrier(ctx)
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
