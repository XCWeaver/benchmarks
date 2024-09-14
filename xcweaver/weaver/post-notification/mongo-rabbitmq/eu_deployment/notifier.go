package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ServiceWeaver/weaver"
	amqp "github.com/rabbitmq/amqp091-go"
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
	weaver.Implements[Notifier]
	weaver.WithConfig[notifierOptions]
	conn *amqp.Connection
}

type notifierOptions struct {
	RabbitMQAddr     string `toml:"rabbitmq_address"`
	RabbitMQPort     string `toml:"rabbitmq_port"`
	RabbitMQUser     string `toml:"rabbitmq_user"`
	RabbitMQPassword string `toml:"rabbitmq_password"`
}

func (n *notifier) Init(ctx context.Context) error {
	logger := n.Logger(ctx)

	var err error
	n.conn, err = amqp.Dial("amqp://" + n.Config().RabbitMQUser + ":" + n.Config().RabbitMQPassword + "@" + n.Config().RabbitMQAddr + ":" + n.Config().RabbitMQPort + "/")
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", "msg", err.Error())
		return err
	}

	logger.Info("notifier service at eu running!", "host", n.Config().RabbitMQAddr, "port", n.Config().RabbitMQPort, "user", n.Config().RabbitMQUser, "password", n.Config().RabbitMQPassword)

	return nil
}

func (n *notifier) Notify(ctx context.Context, postId Post_id_obj, userId int) error {
	logger := n.Logger(ctx)

	notificationStartTimeMs := time.Now().UnixMilli()
	channel, err := n.conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel", "msg", err.Error())
		return err
	}
	defer channel.Close()

	err = channel.ExchangeDeclare("notifier", "topic", false, false, false, false, nil)
	if err != nil {
		logger.Error("error declaring exchange for rabbitmq", "msg", err.Error())
		// errors close the channel so we force the service to restart
		panic(err)
	}

	message := Message{postId, userId, notificationStartTimeMs}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	amqMsg := amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(messageJSON),
	}

	err = channel.PublishWithContext(ctx, "notifier", "notifier-us", false, false, amqMsg)
	if err != nil {
		logger.Error("error publishing to queue", "routing_key", "notifier-us", "err", err.Error())
	}
	notificationsSent.Inc()
	logger.Debug("nofitication successfully written", "postId", postId.PostId, "userId", userId)

	return nil
}
