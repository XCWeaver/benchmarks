package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/XCWeaver/xcweaver"
)

type Message struct {
	PostId                  Post_id_obj
	UserId                  int
	NotificationStartTimeMs int64
}

// Notifier component.
type Notifier interface {
	ReadNotification(ctx context.Context) error
}

// Implementation of the Notifier component.
type notifier struct {
	xcweaver.Implements[Notifier]
	follower_Notify xcweaver.Ref[FollowerNotify]
	clientRabbitMQ  xcweaver.Antipode
}

func (n *notifier) Init(ctx context.Context) error {
	logger := n.Logger(ctx)
	logger.Info("notifier service at us running!")

	return nil
}

func (n *notifier) processMessage(ctx context.Context, d xcweaver.AntipodeObject) {
	logger := n.Logger(ctx)
	var message Message

	// Unmarshal the JSON message into the struct
	err := json.Unmarshal([]byte(d.Version), &message)
	if err != nil {
		logger.Error("Failed to unmarshal JSON", "msg", err.Error())
		return
	}

	postId := message.PostId
	userId := message.UserId
	notificationStartTimeMs := message.NotificationStartTimeMs
	queueDurationMs.Put(float64(time.Now().UnixMilli() - notificationStartTimeMs))
	logger.Debug("New notification received", "postId", postId, "userId", userId)
	notificationsReceived.Inc()

	ctx, err = xcweaver.Transfer(ctx, d.Lineage)
	if err != nil {
		logger.Error("Error transfering the new lineage to context!", "msg", err.Error())
		return
	}

	err = n.follower_Notify.Get().Follower_Notify(ctx, postId, userId)
	if err != nil {
		return
	}
}

func (n *notifier) ReadNotification(ctx context.Context) error {

	var forever chan struct{}
	var stop chan struct{}

	go func() {
		defer close(stop)

		msgs, err := n.clientRabbitMQ.Consume(ctx, "notifier", "notifications-us", stop)
		if err != nil {
			return
		}

		for d := range msgs {
			go n.processMessage(ctx, d)
		}
	}()

	<-forever
	return nil
}
