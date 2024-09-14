package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type CancelService interface {
	Calculate(ctx context.Context, orderId, token string) (float32, error)
	CancelTicket(ctx context.Context, orderId, loginId, token string) (string, error)
}

type cancelService struct {
	weaver.Implements[CancelService]
	orderService         weaver.Ref[OrderService]
	orderOtherService    weaver.Ref[OrderOtherService]
	userService          weaver.Ref[UserService]
	notificationService  weaver.Ref[NotificationService]
	insidePaymentService weaver.Ref[InsidePaymentService]
	roles                []string
}

func (csi *cancelService) Calculate(ctx context.Context, orderId, token string) (float32, error) {
	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return 0.0, err
	}

	var order util.Order
	order, err = csi.orderService.Get().GetOrderById(ctx, orderId, token)
	if err != nil {
		order, err = csi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
	}

	if err != nil {
		return 0.0, err
	}

	if util.OrderStatus(order.Status) == util.NotPaid {
		return 0.0, errors.New("Nothing to refund")
	} else if util.OrderStatus(order.Status) == util.Paid {
		nowDate := time.Now()
		trDate, _ := time.Parse(time.ANSIC, order.TravelDate)

		if nowDate.After(trDate) {
			return 0.0, nil
		} else {
			price := order.Price * 0.8
			return price, nil
		}
	} else {
		return 0.0, errors.New("Refund not permitted")
	}
}

func (csi *cancelService) CancelTicket(ctx context.Context, orderId, loginId, token string) (string, error) {
	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return "", err
	}

	var order util.Order
	order, err = csi.orderService.Get().GetOrderById(ctx, orderId, token)
	if err != nil {
		order, err = csi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
	}

	if err != nil {
		return "", err
	}

	status := util.OrderStatus(order.Status)
	if status != util.Paid && status != util.NotPaid && status != util.Change {
		return "", errors.New("Cancelation not permitted.")
	}

	nowDate := time.Now()
	trDate, _ := time.Parse(time.ANSIC, order.TravelDate)

	var refund string
	if nowDate.After(trDate) {
		refund = "0"
	} else {
		refund = fmt.Sprintf("%f", order.Price*0.8)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var err1, err2 error

	var user util.User
	go func() {
		defer wg.Done()
		_, err1 = csi.insidePaymentService.Get().DrawBack(ctx, loginId, refund, token)
	}()
	go func() {
		defer wg.Done()
		user, err2 = csi.userService.Get().GetUserById(ctx, order.AccountId, token)
	}()
	wg.Wait()
	if err1 != nil {
		return "", err1
	}
	if err2 != nil {
		return "", err2
	}

	err = csi.notificationService.Get().OrderCancelSuccess(ctx, util.NotificationInfo{
		Email:         user.Email,
		OrderNumber:   order.Id,
		Username:      user.Username,
		StartingPlace: order.From,
		EndPlace:      order.To,
		StartingTime:  order.TravelDate,
		SeatClass:     util.SeatClass(order.SeatClass).String(),
		Price:         fmt.Sprintf("%f", order.Price),
	}, token)

	if err != nil {
		return "", nil
	}

	return "Cancelation successful", nil
}
