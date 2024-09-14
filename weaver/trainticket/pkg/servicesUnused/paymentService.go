package services

import (
	"context"
	"errors"
	"log"

	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PaymentService interface {
	//Query(ctx context.Context, token string) ([]model.Payment, error)
	Pay(ctx context.Context, orderId, price, userId, token string) (model.Payment, error)
	//AddMoney(ctx context.Context, userId, price, token string) (model.Payment, error)
}

/*type paymentServiceOptions struct {
	MongoAddr string `toml:"mongodb_address"`
	MongoPort string `toml:"mongodb_port"`
}*/

type paymentService struct {
	weaver.Implements[PaymentService]
	//weaver.WithConfig[paymentServiceOptions]
	clientOptions *options.ClientOptions
	roles         []string
}

func (psi *paymentService) Init(ctx context.Context) error {
	//logger := psi.Logger(ctx)

	//psi.clientOptions = options.Client().ApplyURI("mongodb://" + psi.Config().MongoAddr + ":" + psi.Config().MongoPort + "/?directConnection=true")

	psi.roles = append(psi.roles, "role1")
	psi.roles = append(psi.roles, "role2")
	psi.roles = append(psi.roles, "role3")

	//logger.Info("payment service running!", "mongodb_addr", psi.Config().MongoAddr, "mongodb_port", psi.Config().MongoPort)
	return nil
}

/*func (psi *orderService) Query(ctx context.Context, token string) ([]model.Payment, error) {
	err := util.Authenticate(token, psi.roles...)
	if err != nil {
		return nil, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("payment")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var payments []model.Payment
	err = result.All(&payments)
	if err != nil {
		return nil, err
	}

	return payments, nil
}*/

func (psi *paymentService) Pay(ctx context.Context, orderId, price, userId, token string) (model.Payment, error) {
	logger := psi.Logger(ctx)
	logger.Info("entering Pay", "orderId", orderId)

	err := util.Authenticate(token, psi.roles...)
	if err != nil {
		return model.Payment{}, err
	}

	client, err := mongo.Connect(ctx, psi.clientOptions)
	if err != nil {
		return model.Payment{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("payment")
	filter := bson.D{{"orderId", orderId}}

	res := collection.FindOne(ctx, filter)
	if res.Err() == nil {
		return model.Payment{}, errors.New("Payment already exists for this order!")
	} else if res.Err() != mongo.ErrNoDocuments && res.Err() != nil {
		return model.Payment{}, res.Err()
	}

	newPayment := model.Payment{
		Id:      uuid.New().String(),
		UserId:  userId,
		OrderId: orderId,
		Price:   price,
	}

	result, err := collection.InsertOne(ctx, newPayment)
	if err != nil {
		return model.Payment{}, err
	}
	logger.Debug("inserted payment", "objectid", result.InsertedID)

	return newPayment, nil
}

// ! This func is deprecated
/*func (psi *orderService) AddMoney(ctx context.Context, userId, price, token string) (model.Payment, error) {
	err := util.Authenticate(token, psi.roles...)
	if err != nil {
		return model.Payment{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("payment")

	moneyToAdd := model.Payment{
		Id:     uuid.New().String(),
		UserId: userId,
		Price:  price,
	}

	err = collection.InsertOne(moneyToAdd)
	if err != nil {
		return model.Payment{}, err
	}

	return moneyToAdd, nil
}*/
