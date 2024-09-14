package services

import (
	"context"
	"errors"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type ExecuteService interface {
	ExecuteTicket(ctx context.Context, orderId, token string) (string, error)
	CollectTicket(ctx context.Context, orderId, token string) (string, error)
}

type executeService struct {
	weaver.Implements[ExecuteService]
	orderService      weaver.Ref[OrderService]
	orderOtherService weaver.Ref[OrderOtherService]
	roles             []string
}

func (esi *executeService) ExecuteTicket(ctx context.Context, orderId, token string) (string, error) {
	err := util.Authenticate(token, esi.roles...)
	if err != nil {
		return "", err
	}

	var order util.Order
	first := true
	order, err = esi.orderService.Get().GetOrderById(ctx, orderId, token)
	if err == nil {

		order, err = esi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
		if err != nil {
			return "", err
		}
		first = false
	}

	if order.Status != uint16(util.Paid) && order.Status != uint16(util.Change) {
		return "", errors.New("util.Order cannot be collected!")
	}

	if first {
		_, err = esi.orderService.Get().ModifyOrder(ctx, orderId, uint16(util.Collected), token)
	} else {
		_, err = esi.orderOtherService.Get().ModifyOrder(ctx, orderId, uint16(util.Collected), token)
	}

	if err != nil {
		return "", err
	}

	return "util.Order collected successfully", nil
}

func (esi *executeService) CollectTicket(ctx context.Context, orderId, loginId, token string) (string, error) {
	err := util.Authenticate(token, esi.roles...)
	if err != nil {
		return "", err
	}

	first := true
	_, err = esi.orderService.Get().GetOrderById(ctx, orderId, token)
	if err == nil {

		_, err = esi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
		if err != nil {
			return "", err
		}
		first = false
	}

	if first {
		_, err = esi.orderService.Get().ModifyOrder(ctx, orderId, uint16(util.Used), token)
	} else {
		_, err = esi.orderOtherService.Get().ModifyOrder(ctx, orderId, uint16(util.Used), token)
	}

	if err != nil {
		return "", err
	}

	return "util.Order executed successfully", nil
}
