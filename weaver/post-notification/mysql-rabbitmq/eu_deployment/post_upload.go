package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

// PostUpload component.
type PostUpload interface {
	Post(context.Context, string, int) (string, error)
}

// Implementation of the PostUpload component.
type postUpload struct {
	weaver.Implements[PostUpload]
	post_storage weaver.Ref[PostStorageEu]
	notifier     weaver.Ref[Notifier]
}

func (p *postUpload) Init(ctx context.Context) error {
	logger := p.Logger(ctx)
	logger.Info("post upload service at eu running!")

	return nil
}

// Forwards the post data to Post_storage component and then sends the post id to
// the Notifier component
func (p *postUpload) Post(ctx context.Context, post string, userId int) (string, error) {

	//send post to post_storage
	postId, err := p.post_storage.Get().Post(ctx, post)
	if err != nil {
		return "", err
	}

	//send postID and userId to notifier
	err = p.notifier.Get().Notify(ctx, postId, userId)
	if err != nil {
		return "", err
	}

	return postId.PostId, nil
}
