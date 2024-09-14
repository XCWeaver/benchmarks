package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TiagoMalhadas/xcweaver"
)

type Message struct {
	PostId                  Post_id_obj
	UserId                  int
	NotificationStartTimeMs int64
}

// Server component.
type Notifier interface {
	Notify(context.Context, Post_id_obj, int) error
}

// Implementation of the Notifier component.
type notifier struct {
	xcweaver.Implements[Notifier]
	clientRabbitMQ xcweaver.Antipode
}

func (n *notifier) Init(ctx context.Context) error {
	logger := n.Logger(ctx)
	logger.Info("notifier service at eu running!")

	return nil
}

func (n *notifier) Notify(ctx context.Context, postId Post_id_obj, userId int) error {
	logger := n.Logger(ctx)

	notificationStartTimeMs := time.Now().UnixMilli()
	message := Message{postId, userId, notificationStartTimeMs}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Convert the byte slice to a string
	messageString := string(messageJSON)

	_, err = n.clientRabbitMQ.Write(ctx, "notifier", "notifications-us", messageString)
	if err != nil {
		logger.Error("Error writing notification to queue", "msg", err.Error())
		return err
	}
	notificationsSent.Inc()
	logger.Debug("nofitication successfully written", "postId", postId.PostId, "userId", userId)

	return nil
}
