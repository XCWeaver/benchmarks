package services

import (
	"context"
	"errors"
	"fmt"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type PaymentService interface {
	Query(ctx context.Context, token string) ([]util.Payment, error)
	Pay(ctx context.Context, orderId, price, userId, token string) (util.Payment, error)
	AddMoney(ctx context.Context, userId, price, token string) (util.Payment, error)
}

type paymentService struct {
	weaver.Implements[PaymentService]
	//Mongo
	db    components.NoSQLDatabase
	roles []string
}

func (psi *orderService) Query(ctx context.Context, token string) ([]util.Payment, error) {
	err := util.Authenticate(token, psi.roles...)
	if err != nil {
		return nil, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("payment")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var payments []util.Payment
	err = result.All(&payments)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (psi *orderService) Pay(ctx context.Context, orderId, price, userId, token string) (util.Payment, error) {
	err := util.Authenticate(token, psi.roles...)
	if err != nil {
		return util.Payment{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("payment")
	query := fmt.Sprintf(`{"OrderId": %s }`, orderId)

	res, err := collection.FindOne(query)
	if err == nil {
		var oldPayment util.Payment
		res.Decode(&oldPayment)

		if oldPayment.Id != "" {
			return util.Payment{}, errors.New("util.Payment already exists for this order!")
		}
	}

	newPayment := util.Payment{
		Id:      uuid.New().String(),
		UserId:  userId,
		OrderId: orderId,
		Price:   price,
	}

	err = collection.InsertOne(newPayment)
	if err != nil {
		return util.Payment{}, err
	}

	return newPayment, nil
}

// ! This func is deprecated
func (psi *orderService) AddMoney(ctx context.Context, userId, price, token string) (util.Payment, error) {
	err := util.Authenticate(token, psi.roles...)
	if err != nil {
		return util.Payment{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("payment")

	moneyToAdd := util.Payment{
		Id:     uuid.New().String(),
		UserId: userId,
		Price:  price,
	}

	err = collection.InsertOne(moneyToAdd)
	if err != nil {
		return util.Payment{}, err
	}

	return moneyToAdd, nil
}
