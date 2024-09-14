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

// Notifier component.
type Notifier interface {
}

// Implementation of the Notifier component.
type notifier struct {
	weaver.Implements[Notifier]
	follower_Notify weaver.Ref[Follower_Notify]
	weaver.WithConfig[notifierOptions]
}

type notifierOptions struct {
	RabbitMQAddr     string `toml:"rabbitmq_address"`
	RabbitMQPort     string `toml:"rabbitmq_port"`
	RabbitMQUser     string `toml:"rabbitmq_user"`
	RabbitMQPassword string `toml:"rabbitmq_password"`
}

func (n *notifier) Init(ctx context.Context) error {
	logger := n.Logger(ctx)

	conn, err := amqp.Dial("amqp://" + n.Config().RabbitMQUser + ":" + n.Config().RabbitMQPassword + "@" + n.Config().RabbitMQAddr + ":" + n.Config().RabbitMQPort + "/")
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", "msg", err.Error())
		return err
	}

	logger.Info("notifier service at us running!", "host", n.Config().RabbitMQAddr, "port", n.Config().RabbitMQPort, "user", n.Config().RabbitMQUser, "password", n.Config().RabbitMQPassword)

	err = n.readNotification(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

func (n *notifier) processMessage(ctx context.Context, msg amqp.Delivery) {
	logger := n.Logger(ctx)
	var message Message

	// Unmarshal the JSON message into the struct
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		logger.Error("Failed to unmarshal JSON", "msg", err.Error())
		return
	}

	postId := message.PostId
	userId := message.UserId
	notificationStartTimeMs := message.NotificationStartTimeMs
	queueDurationMs.Put(float64(time.Now().UnixMilli() - notificationStartTimeMs))
	notificationsReceived.Inc()

	n.follower_Notify.Get().Follower_Notify(ctx, postId, userId)
}

func (n *notifier) readNotification(ctx context.Context, conn *amqp.Connection) error {
	logger := n.Logger(ctx)

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a channel", "msg", err.Error())
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("notifier", "topic", false, false, false, false, nil)
	if err != nil {
		logger.Error("error declaring exchange for rabbitmq", "msg", err.Error())
		return err
	}
	routingKey := "notifier-us"
	_, err = ch.QueueDeclare(routingKey, false, false, false, false, nil)
	if err != nil {
		logger.Error("error declaring queue for rabbitmq", "msg", err.Error())
		return err
	}

	err = ch.QueueBind(routingKey, routingKey, "notifier", false, nil)
	if err != nil {
		logger.Error("error binding queue for rabbitmq", "msg", err.Error())
		return err
	}

	msgs, err := ch.Consume(routingKey, "", true, false, false, false, nil)
	if err != nil {
		logger.Error("error consuming queue", "msg", err.Error())
		return err
	}

	for d := range msgs {
		go n.processMessage(ctx, d)
	}
	return nil
}
