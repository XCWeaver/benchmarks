package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

// Follower_Notify component.
type Follower_Notify interface {
	Follower_Notify(context.Context, Post_id_obj, int) error
}

// Implementation of the Follower_Notify component.
type follower_Notify struct {
	weaver.Implements[Follower_Notify]
	post_storage weaver.Ref[Post_storage_america]
}

func (f *follower_Notify) Init(ctx context.Context) error {
	logger := f.Logger(ctx)

	logger.Info("follower notify service running!")
	return nil
}

func (f *follower_Notify) Follower_Notify(ctx context.Context, postId Post_id_obj, userId int) error {
	logger := f.Logger(ctx)

	post, err := f.post_storage.Get().GetPost(ctx, postId)
	if err != nil {
		return err
	}

	logger.Debug("Post read successfully!", "post", post)

	return nil
}
