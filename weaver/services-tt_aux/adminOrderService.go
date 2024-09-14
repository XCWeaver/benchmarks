package services

import (
	"context"
	"errors"
	"sync"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

// ! TODO
type AdminOrderService interface {
	GetAllOrders(ctx context.Context, token string) ([]util.Order, error)
	DeleteOrder(ctx context.Context, orderId, token string) (string, error)
	UpdateOrder(ctx context.Context, order util.Order, token string) (util.Order, error)
	AddOrder(ctx context.Context, order util.Order, token string) (util.Order, error)
}

type adminOrderService struct {
	weaver.Implements[AdminOrderService]
	orderService      weaver.Ref[OrderService]
	orderOtherService weaver.Ref[OrderOtherService]
	roles             []string
}

func (aosi *adminOrderService) GetAllOrders(ctx context.Context, token string) ([]util.Order, error) {
	err := util.Authenticate(token, aosi.roles...)
	if err != nil {
		return nil, err
	}

	var err1, err2 error
	var ordersFirstBatch, ordersSecondBatch []util.Order
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ordersFirstBatch, err1 = aosi.orderService.Get().FindAllOrder(ctx, token)
	}()

	go func() {
		defer wg.Done()
		ordersSecondBatch, err2 = aosi.orderOtherService.Get().FindAllOrder(ctx, token)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	orders := append(ordersFirstBatch, ordersSecondBatch...)

	if len(orders) == 0 {
		return nil, errors.New("No orders found")
	}

	return orders, nil
}

func (aosi *adminOrderService) DeleteOrder(ctx context.Context, orderId, trainNumber, token string) (string, error) {
	err := util.Authenticate(token, aosi.roles...)
	if err != nil {
		return "", err
	}

	var msg string
	if trainNumber[0:1] == "D" || trainNumber[0:1] == "G" {
		msg, err = aosi.orderService.Get().DeleteOrder(ctx, orderId, token)

	} else {
		msg, err = aosi.orderOtherService.Get().DeleteOrder(ctx, orderId, token)
	}

	if err != nil {
		return "", err
	}

	return msg, nil
}

func (aosi *adminOrderService) UpdateOrder(ctx context.Context, order util.Order, token string) (util.Order, error) {
	err := util.Authenticate(token, aosi.roles...)
	if err != nil {
		return util.Order{}, err
	}

	if order.TrainNumber[0:1] == "D" || order.TrainNumber[0:1] == "G" {
		order, err = aosi.orderService.Get().UpdateOrder(ctx, order, token)
	} else {
		order, err = aosi.orderOtherService.Get().UpdateOrder(ctx, order, token)
	}

	if err != nil {
		return util.Order{}, err
	}

	return order, nil
}

func (aosi *adminOrderService) AddOrder(ctx context.Context, order util.Order, token string) (util.Order, error) {
	err := util.Authenticate(token, aosi.roles...)
	if err != nil {
		return util.Order{}, err
	}

	if order.TrainNumber[0:1] == "D" || order.TrainNumber[0:1] == "G" {
		order, err = aosi.orderService.Get().AddCreateNewOrder(ctx, order, token)
	} else {
		order, err = aosi.orderOtherService.Get().AddCreateNewOrder(ctx, order, token)
	}

	if err != nil {
		return util.Order{}, err
	}

	return order, nil
}
