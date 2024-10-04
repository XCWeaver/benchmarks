package main

import (
	"context"

	"github.com/XCWeaver/xcweaver"
)

// PostUpload component.
type PostUpload interface {
	Post(context.Context, string, int) (string, error)
}

// Implementation of the PostUpload component.
type postUpload struct {
	xcweaver.Implements[PostUpload]
	post_storage xcweaver.Ref[PostStorageEu]
	notifier     xcweaver.Ref[Notifier]
}

func (p *postUpload) Init(ctx context.Context) error {
	logger := p.Logger(ctx)
	logger.Info("post upload service at eu running!")

	return nil
}

// Forwards the post data to Post_storage component and then sends the post id to
// the Notifier component
func (p *postUpload) Post(ctx context.Context, post string, userId int) (string, error) {
	logger := p.Logger(ctx)

	//send post to post_storage
	lineage, postId, err := p.post_storage.Get().Post(ctx, post)
	if err != nil {
		return "", err
	}

	ctx, err = xcweaver.Transfer(ctx, lineage)
	if err != nil {
		logger.Error("Error transfering lineage to context!", "msg", err.Error())
		return "", err
	}

	//send postID and userId to notifier
	err = p.notifier.Get().Notify(ctx, postId, userId)
	if err != nil {
		return "", err
	}

	return postId.PostId, nil
}
