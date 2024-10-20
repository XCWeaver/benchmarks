package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	tt_metrics "trainticket/pkg/metrics"
	"trainticket/pkg/model"

	"github.com/ServiceWeaver/weaver"
)

type CancelService interface {
	Calculate(ctx context.Context, orderId, token string) (string, error)
	CancelTicket(ctx context.Context, orderId, loginId, token string) (string, error)
	GetConsistencyWindowValues(ctx context.Context) ([]float64, error)
	GetInconsistencies(ctx context.Context) (int, error)
	Reset(ctx context.Context) error
}

type cancelService struct {
	weaver.Implements[CancelService]
	orderService      weaver.Ref[OrderService]
	orderOtherService weaver.Ref[OrderOtherService]
	//userService       weaver.Ref[UserService]
	//notificationService     weaver.Ref[NotificationService]
	insidePaymentService    weaver.Ref[InsidePaymentService]
	roles                   []string
	mu                      sync.Mutex
	consistencyWindowValues []float64
	muInc                   sync.Mutex
	inconsistencies         int
}

func (csi *cancelService) Init(ctx context.Context) error {
	logger := csi.Logger(ctx)

	csi.roles = append(csi.roles, "role1")
	csi.roles = append(csi.roles, "role2")
	csi.roles = append(csi.roles, "role3")

	csi.inconsistencies = 0

	logger.Info("cancel service running!")
	return nil
}

func (csi *cancelService) Calculate(ctx context.Context, orderId, token string) (string, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering Calculate", "orderId", orderId)

	/*err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return "", err
	}*/

	var order model.Order
	order, err := csi.orderService.Get().GetOrderById(ctx, orderId, token)
	if err == nil {
		if model.OrderStatus(order.Status) == model.NotPaid || model.OrderStatus(order.Status) == model.Paid {
			if model.OrderStatus(order.Status) == model.NotPaid {
				logger.Debug("Success. Refoud 0")
				return "0.0", nil
			} else {
				logger.Debug("[Cancel Order][Refund Price] From Order Service.Paid.")
				return calculateRefund(order), nil
			}
		} else {
			return "", errors.New("Order Status Cancel Not Permitted, Refund error")
		}
	} else {
		order, err = csi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
		if err != nil {
			return "", err
		}
		if model.OrderStatus(order.Status) == model.NotPaid || model.OrderStatus(order.Status) == model.Paid {
			if model.OrderStatus(order.Status) == model.NotPaid {
				logger.Debug("Success. Refoud 0")
				return "0.0", nil
			} else {
				logger.Debug("[Cancel Order][Refund Price] From Order Service.Paid.")
				return calculateRefund(order), nil
			}
		} else {
			return "", errors.New("Order Status Cancel Not Permitted, Refund error")
		}
	}
}

func (csi *cancelService) CancelTicket(ctx context.Context, orderId, loginId, token string) (string, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering CancelTicket", "orderId", orderId)

	/*err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return "", err
	}*/
	startTimeMs := time.Now().UnixMilli()

	order, err := csi.orderService.Get().GetOrderById(ctx, orderId, token)
	if err == nil {
		status := model.OrderStatus(order.Status)
		if status == model.Paid || status == model.NotPaid || status == model.Change || status == model.Cancel {
			/*********************** Fault Reproduction - Error Process Seq *************************/
			// 1. cancel order
			_, err = csi.orderService.Get().ModifyOrder(ctx, orderId, 4, token)
			if err != nil {
				return "", err
			}

			// 2. drawback money
			refund := calculateRefund(order)
			var wg sync.WaitGroup
			wg.Add(1)
			var err1 error

			alive := make([]int32, 1)
			atomic.StoreInt32(&alive[0], 1)
			go func() {
				defer wg.Done()
				defer atomic.StoreInt32(&alive[0], 0)
				_, err1 = csi.insidePaymentService.Get().DrawBack(ctx, loginId, refund, token)
			}()

			//wg.Wait()

			consistencyWindowMs := float64(time.Now().UnixMilli() - startTimeMs)
			tt_metrics.ConsistencyWindow.Put(consistencyWindowMs)
			csi.mu.Lock()
			csi.consistencyWindowValues = append(csi.consistencyWindowValues, consistencyWindowMs)
			csi.mu.Unlock()

			// 3. get results
			if alive[0] == 1 {
				logger.Error("Assynchronous call not yet finished!")
				tt_metrics.Inconsistencies.Inc()
				csi.muInc.Lock()
				csi.inconsistencies += 1
				csi.muInc.Unlock()
				return "Cancelation failed", nil
			}
			if err1 != nil {
				return "", err1
			}
		}
	} else {
		order, err = csi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
		if err != nil {
			return "", err
		}
		/*********************** Fault Reproduction - Error Process Seq *************************/
		// 1. cancel order
		_, err = csi.orderOtherService.Get().ModifyOrder(ctx, orderId, 4, token)
		if err != nil {
			return "", err
		}

		// 2. drawback money
		refund := calculateRefund(order)
		var wg sync.WaitGroup
		wg.Add(1)
		var err1 error

		alive := make([]int32, 1)
		atomic.StoreInt32(&alive[0], 1)
		go func() {
			defer wg.Done()
			defer atomic.StoreInt32(&alive[0], 0)
			_, err1 = csi.insidePaymentService.Get().DrawBack(ctx, loginId, refund, token)
		}()

		//wg.Wait()

		consistencyWindowMs := float64(time.Now().UnixMilli() - startTimeMs)
		tt_metrics.ConsistencyWindow.Put(consistencyWindowMs)
		csi.mu.Lock()
		csi.consistencyWindowValues = append(csi.consistencyWindowValues, consistencyWindowMs)
		csi.mu.Unlock()
		// 3. get results
		if alive[0] == 1 {
			logger.Error("Assynchronous call not yet finished!")
			tt_metrics.Inconsistencies.Inc()
			csi.muInc.Lock()
			csi.inconsistencies += 1
			csi.muInc.Unlock()
			return "Cancelation failed", nil
		}
		if err1 != nil {
			return "", err1
		}
	}

	return "Cancelation successful", nil
}

func (csi *cancelService) GetConsistencyWindowValues(ctx context.Context) ([]float64, error) {
	csi.mu.Lock()
	values := csi.consistencyWindowValues
	csi.mu.Unlock()
	return values, nil
}

func (csi *cancelService) GetInconsistencies(ctx context.Context) (int, error) {
	csi.muInc.Lock()
	inconsistencies := csi.inconsistencies
	csi.muInc.Unlock()
	return inconsistencies, nil
}

func (csi *cancelService) Reset(ctx context.Context) error {
	csi.inconsistencies = 0
	csi.consistencyWindowValues = []float64{}
	return nil
}

func calculateRefund(order model.Order) string {
	nowDate := time.Now()
	trDate, _ := time.Parse(time.ANSIC, order.TravelDate)

	if nowDate.After(trDate) {
		return "0"
	} else {
		return fmt.Sprintf("%f", order.Price*0.8)
	}
}
