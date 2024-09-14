package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type InsidePaymentService interface {
	Pay(ctx context.Context, tripId, userId, orderId, token string) (string, error)
	CreateAccount(ctx context.Context, money, userId, token string) (string, error)
	AddMoney(ctx context.Context, userId, money, token string) (string, error)
	QueryPayment(ctx context.Context, token string) ([]util.Payment, error)
	QueryAccount(ctx context.Context, token string) ([]util.Balances, error)
	DrawBack(ctx context.Context, userId, money, token string) (string, error)
	PayDifference(ctx context.Context, orderId, userId, price, token string) (string, error)
	QueryAddMoney(ctx context.Context, token string) ([]util.AddMoney, error)
}

type insidePaymentService struct {
	weaver.Implements[InsidePaymentService]
	//Mongo
	db                components.NoSQLDatabase
	orderService      weaver.Ref[OrderService]
	orderOtherService weaver.Ref[OrderOtherService]
	paymentService    weaver.Ref[PaymentService]
	roles             []string
}

func (ipsi *insidePaymentService) Pay(ctx context.Context, tripId, userId, orderId, token string) (string, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return "", err
	}

	var order util.Order
	if tripId[0:1] == "G" || tripId[0:1] == "D" {
		order, err = ipsi.orderService.Get().GetOrderById(ctx, orderId, token)
	} else {
		order, err = ipsi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
	}

	newPayment := util.Payment{
		Id:      uuid.New().String(),
		OrderId: orderId,
		UserId:  userId,
		Price:   fmt.Sprintf("%f", order.Price),
	}

	query := fmt.Sprintf(`{"UserId": %s }`, userId)
	collection := ipsi.db.GetDatabase("ts").GetCollection("payment")
	res, err := collection.FindMany(query)
	if err != nil {
		return "", err
	}
	var payments []util.Payment
	res.All(&payments)

	totalExpand := order.Price
	for _, p := range payments {
		price, _ := strconv.ParseFloat(p.Price, 32)
		totalExpand += float32(price)
	}

	amCollection := ipsi.db.GetDatabase("ts").GetCollection("addMoney")
	res, err = amCollection.FindMany(query)
	if err != nil {
		return "", err
	}
	var addMoney []util.AddMoney
	res.Decode(&addMoney)

	totalMoney := float32(0.0)
	for _, am := range addMoney {
		money, _ := strconv.ParseFloat(am.Money, 32)
		totalMoney += float32(money)
	}

	if totalExpand > totalMoney {
		_, err = ipsi.paymentService.Get().Pay(ctx, orderId, fmt.Sprintf("%f", order.Price), userId, token)
		if err != nil {
			return "", err
		}
		newPayment.Type = util.OutsidePayment.String()
	} else {
		newPayment.Type = util.NormalPayment.String()
	}

	err = collection.InsertOne(newPayment)
	if err != nil {
		return "", err
	}

	if tripId[0:1] == "G" || tripId[0:1] == "D" {
		_, err = ipsi.orderService.Get().ModifyOrder(ctx, orderId, uint16(util.Paid), token)
	} else {
		_, err = ipsi.orderOtherService.Get().ModifyOrder(ctx, orderId, uint16(util.Paid), token)
	}

	if err != nil {
		return "", err
	}

	return "util.Payment successful", nil
}

func (ipsi *insidePaymentService) CreateAccount(ctx context.Context, money, userId, token string) (string, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return "", err
	}

	collection := ipsi.db.GetDatabase("ts").GetCollection("addMoney")
	query := fmt.Sprintf(`{"UserId": %s }`, userId)

	res, err := collection.FindOne(query)

	if err == nil {
		var am util.AddMoney
		err = res.Decode(&am)
		if err == nil {
			return "", errors.New("Account already exists for this user.")
		}
	}

	err = collection.InsertOne(util.AddMoney{
		Id:     uuid.New().String(),
		Money:  money,
		UserId: userId,
		Type:   util.AddMoneyType.String(),
	})
	if err != nil {
		return "", err
	}

	return "Account created successfully.", nil

}

func (ipsi *insidePaymentService) AddMoney(ctx context.Context, userId, money, token string) (string, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return "", err
	}

	collection := ipsi.db.GetDatabase("ts").GetCollection("addMoney")

	query := fmt.Sprintf(`{"UserId": %s }`, userId)

	res, err := collection.FindOne(query)
	if err != nil {
		return "", err
	}

	var account util.AddMoney
	res.Decode(&account)

	uQuery := fmt.Sprintf(`{"Id": %s }`, account.Id)
	update := fmt.Sprintf(`{"$set": {Money: %s, Type: %s}}`, money, util.AddMoneyType.String())
	err = collection.UpdateOne(uQuery, update)
	if err != nil {
		return "", err
	}

	return "Money added successfully", nil
}

func (ipsi *insidePaymentService) QueryPayment(ctx context.Context, token string) ([]util.Payment, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return nil, err
	}

	collection := ipsi.db.GetDatabase("ts").GetCollection("payment")

	res, err := collection.FindMany("")
	if err != nil {
		return nil, err
	}

	var payments []util.Payment
	err = res.All(&payments)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (ipsi *insidePaymentService) QueryAccount(ctx context.Context, token string) ([]util.Balances, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return nil, err
	}

	collection := ipsi.db.GetDatabase("ts").GetCollection("addMoney")

	res, err := collection.FindMany("")
	if err != nil {
		return nil, err
	}

	var accounts []util.AddMoney
	err = res.All(&accounts)
	if err != nil {
		return nil, err
	}

	moneyMap := make(map[string]float32)

	for _, acc := range accounts {
		toAdd, _ := strconv.ParseFloat(acc.Money, 32)

		if _, ok := moneyMap[acc.UserId]; ok {
			moneyMap[acc.UserId] += float32(toAdd)
		} else {
			moneyMap[acc.UserId] = float32(toAdd)
		}
	}

	paymentsCollection := ipsi.db.GetDatabase("ts").GetCollection("payment")

	var resultBalances []util.Balances

	var totalExpand float32
	for userId, _ := range moneyMap {
		query := fmt.Sprintf(`{"UserId": %s }`, userId)
		userRes, err := paymentsCollection.FindOne(query)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var paymentList []util.Payment
		err = userRes.All(&paymentList)
		if err != nil {
			continue
		}
		totalExpand = 0.0
		for _, p := range paymentList {
			price, _ := strconv.ParseFloat(p.Price, 32)
			totalExpand += float32(price)
		}

		resultBalances = append(resultBalances, util.Balances{
			UserId:  userId,
			Balance: totalExpand,
		})

	}

	return resultBalances, nil

}

func (ipsi *insidePaymentService) DrawBack(ctx context.Context, userId, money, token string) (string, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return "", err
	}

	collection := ipsi.db.GetDatabase("ts").GetCollection("addMoney")
	query := fmt.Sprintf(`{"UserId": %s }`, userId)

	_, err = collection.FindOne(query)
	if err != nil {
		return "", err
	}

	err = collection.InsertOne(util.AddMoney{
		Id:     uuid.New().String(),
		UserId: userId,
		Money:  money,
		Type:   util.DrawBackMoney.String(),
	})
	if err != nil {
		return "", err
	}

	return "Drawback successful", nil
}

func (ipsi *insidePaymentService) PayDifference(ctx context.Context, orderId, userId, price, token string) (string, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return "", err
	}

	collection := ipsi.db.GetDatabase("ts").GetCollection("payment")
	query := fmt.Sprintf(`{"UserId": %s }`, userId)

	res, err := collection.FindMany(query)
	if err != nil {
		return "", err
	}
	var payments []util.Payment
	err = res.All(&payments)
	if err != nil {
		return "", err
	}

	totalExpand, _ := strconv.ParseFloat(price, 32)

	for _, p := range payments {
		preyes, _ := strconv.ParseFloat(p.Price, 32)
		totalExpand += preyes
	}

	amCollection := ipsi.db.GetDatabase("ts").GetCollection("addMoney")
	res, err = amCollection.FindMany(query)
	if err != nil {
		return "", err
	}
	var accounts []util.AddMoney
	err = res.All(&accounts)
	if err != nil {
		return "", err
	}

	totalMoney := float32(0.0)
	for _, a := range accounts {
		money, _ := strconv.ParseFloat(a.Money, 32)
		totalMoney += float32(money)
	}

	if float32(totalExpand) > totalMoney {
		ipsi.paymentService.Get().Pay(ctx, orderId, userId, price, token)
	}

	newPayment := util.Payment{
		Id:      uuid.New().String(),
		UserId:  userId,
		OrderId: orderId,
		Price:   price,
		Type:    util.ExternalAndDifferencePayment.String(),
	}

	err = collection.InsertOne(newPayment)
	if err != nil {
		return "", err
	}

	return "Difference payment successful.", nil
}

func (ipsi *insidePaymentService) QueryAddMoney(ctx context.Context, token string) ([]util.AddMoney, error) {
	err := util.Authenticate(token, ipsi.roles...)
	if err != nil {
		return nil, err
	}

	collection := ipsi.db.GetDatabase("ts").GetCollection("addMoney")

	res, err := collection.FindMany("")
	if err != nil {
		return nil, err
	}
	var accounts []util.AddMoney

	err = res.All(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
