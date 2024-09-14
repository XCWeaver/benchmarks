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

type RebookService interface {
	PayDifference(ctx context.Context, info util.RebookInfo, token string) (util.Order, error)
	Rebook(ctx context.Context, info util.RebookInfo, token string) (util.Order, error)
}

type rebookService struct {
	weaver.Implements[RebookService]
	orderService         weaver.Ref[OrderService]
	orderOtherService    weaver.Ref[OrderOtherService]
	insidePaymentService weaver.Ref[InsidePaymentService]
	seatService          weaver.Ref[SeatService]
	stationService       weaver.Ref[StationService]
	travelService        weaver.Ref[TravelService]
	travel2Service       weaver.Ref[Travel2Service]
	roles                []string
}

func (rsi *rebookService) IsTripGD(tripId string) bool {
	return tripId[:1] == "G" || tripId[:1] == "D"
}

func (rsi *rebookService) PayDifference(ctx context.Context, info util.RebookInfo, token string) (util.Order, error) {
	err := util.Authenticate(token, rsi.roles...)
	if err != nil {
		return util.Order{}, err
	}

	var order util.Order
	if rsi.IsTripGD(info.OldTripId) {
		order, err = rsi.orderService.Get().GetOrderById(ctx, info.OrderId, token)
	} else {
		order, err = rsi.orderOtherService.Get().GetOrderById(ctx, info.OrderId, token)
	}

	if err != nil {
		return util.Order{}, err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var err1, err2 error

	var from, to string
	go func() {
		defer wg.Done()
		from, err1 = rsi.stationService.Get().QueryById(ctx, order.From, token)
	}()
	go func() {
		defer wg.Done()
		to, err2 = rsi.stationService.Get().QueryById(ctx, order.To, token)
	}()
	wg.Wait()
	if err1 != nil {
		return util.Order{}, err1
	}
	if err2 != nil {
		return util.Order{}, err2
	}

	var tripResponse util.TripDetails

	if rsi.IsTripGD(info.TripId) {
		_, tripResponse, err = rsi.travelService.Get().GetTripAllDetailInfo(ctx, info.TripId, from, to, info.Date, token)
	} else {
		_, tripResponse, err = rsi.travel2Service.Get().GetTripAllDetailInfo(ctx, info.TripId, from, to, info.Date, token)
	}

	var ticketPrice float32
	if info.SeatType == uint16(util.FirstClass) {
		ticketPrice = tripResponse.PriceForComfortClass
	} else if info.SeatType == uint16(util.SecondClass) {
		ticketPrice = tripResponse.PriceForEconomyClass
	}

	_, err = rsi.insidePaymentService.Get().PayDifference(ctx, info.OrderId, info.LoginId, fmt.Sprintf("%f", ticketPrice-order.Price), token)
	if err != nil {
		return util.Order{}, err
	}

	//* Identical to rebook func part
	_, err = rsi.seatService.Get().Create(ctx, info.SeatType, info.Date, order.TrainNumber, order.From, order.To, token)

	if err != nil {
		return util.Order{}, err
	}

	//! Some details about the `date` parameter are changed here from the original system
	//! The `util.Order` object in the original system kept two separate attributes for date and time respectively.
	//! We need to verify any use cases that might be affected by this. For now, we just create orders
	//! with a "From" datetime equal to <date><trip.startingTime[time]>. We get the time portion only from startingTime.

	oldTripId := order.TrainNumber

	tmpDate, _ := time.Parse(time.ANSIC, info.Date)
	y, m, d := tmpDate.Date()
	startTime, _ := time.Parse(time.ANSIC, tripResponse.StartingTime)
	hr, min, sec := startTime.Clock()
	trDate := fmt.Sprintf("%s %s %s %s:%s:%s %s", tmpDate.Weekday().String()[:3], m.String()[:3], d, hr, min, sec, y)
	bgDate := time.Now()

	order.BoughtDate = bgDate.String()
	order.TravelDate = trDate
	order.TrainNumber = info.TripId
	order.SeatClass = info.SeatType
	order.Status = uint16(util.Change)
	order.Price = ticketPrice

	if (rsi.IsTripGD(oldTripId) && rsi.IsTripGD(info.TripId)) || (!rsi.IsTripGD(oldTripId) && !rsi.IsTripGD(info.TripId)) {

		if rsi.IsTripGD(info.TripId) {
			order, err = rsi.orderService.Get().SaveOrderInfo(ctx, order, token)
		} else {
			order, err = rsi.orderOtherService.Get().SaveOrderInfo(ctx, order, token)
		}
		if err != nil {
			return util.Order{}, err
		}
	} else {

		if rsi.IsTripGD(oldTripId) {
			_, err = rsi.orderService.Get().DeleteOrder(ctx, order.Id, token)
		} else {
			_, err = rsi.orderOtherService.Get().DeleteOrder(ctx, order.Id, token)
		}

		if err != nil {
			return util.Order{}, err
		}

		if rsi.IsTripGD(info.TripId) {
			order, err = rsi.orderService.Get().CreateNewOrder(ctx, order, token)
		} else {
			order, err = rsi.orderOtherService.Get().CreateNewOrder(ctx, order, token)
		}
		if err != nil {
			return util.Order{}, err
		}
	}

	return order, nil
}

func (rsi *rebookService) Rebook(ctx context.Context, info util.RebookInfo, token string) (util.Order, error) {
	err := util.Authenticate(token, rsi.roles...)
	if err != nil {
		return util.Order{}, err
	}

	var order util.Order
	if rsi.IsTripGD(info.OldTripId) {
		order, err = rsi.orderService.Get().GetOrderById(ctx, info.OrderId, token)
	} else {
		order, err = rsi.orderOtherService.Get().GetOrderById(ctx, info.OrderId, token)
	}

	if err != nil {
		return util.Order{}, err
	}

	if order.Status != uint16(util.Paid) {
		return util.Order{}, errors.New("util.Order cannot be rebooked.")
	}

	travelDateTime, _ := time.Parse(time.ANSIC, order.TravelDate)
	diff := time.Now().Sub(travelDateTime)

	if diff.Hours() > 2 {
		return util.Order{}, errors.New("util.Order can be rebooked only up to two hours after travel time.")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var err1, err2 error

	var from, to string
	go func() {
		defer wg.Done()
		from, err1 = rsi.stationService.Get().QueryById(ctx, order.From, token)
	}()
	go func() {
		defer wg.Done()
		to, err2 = rsi.stationService.Get().QueryById(ctx, order.To, token)
	}()
	wg.Wait()
	if err1 != nil {
		return util.Order{}, err1
	}
	if err2 != nil {
		return util.Order{}, err2
	}

	var tripResponse util.TripDetails

	if rsi.IsTripGD(info.TripId) {
		_, tripResponse, err = rsi.travelService.Get().GetTripAllDetailInfo(ctx, info.TripId, from, to, info.Date, token)
	} else {
		_, tripResponse, err = rsi.travel2Service.Get().GetTripAllDetailInfo(ctx, info.TripId, from, to, info.Date, token)
	}

	var ticketPrice float32
	if info.SeatType == uint16(util.FirstClass) {
		ticketPrice = tripResponse.PriceForComfortClass
	} else if info.SeatType == uint16(util.SecondClass) {
		ticketPrice = tripResponse.PriceForEconomyClass
	}

	oldPrice := order.Price

	if oldPrice > ticketPrice {
		difference := fmt.Sprintf("%f", oldPrice-ticketPrice)
		_, err = rsi.insidePaymentService.Get().DrawBack(ctx, info.LoginId, difference, token)
		if err != nil {
			return util.Order{}, err
		}

	} else if oldPrice < ticketPrice {
		return util.Order{}, errors.New(fmt.Sprintf("Please pay difference.%f", ticketPrice-oldPrice))
	}

	//* Identical to rebook func part
	_, err = rsi.seatService.Get().Create(ctx, info.SeatType, info.Date, order.TrainNumber, order.From, order.To, token)

	if err != nil {
		return util.Order{}, err
	}

	//! Some details about the `date` parameter are changed here from the original system
	//! The `util.Order` object in the original system kept two separate attributes for date and time respectively.
	//! We need to verify any use cases that might be affected by this. For now, we just create orders
	//! with a "From" datetime equal to <date><trip.startingTime[time]>. We get the time portion only from startingTime.

	oldTripId := order.TrainNumber

	tmpDate, _ := time.Parse(time.ANSIC, info.Date)
	y, m, d := tmpDate.Date()
	hr, min, sec := tmpDate.Clock()
	trDate := fmt.Sprintf("%s %s %s %s:%s:%s %s", tmpDate.Weekday().String()[:3], m.String()[:3], d, hr, min, sec, y)
	bgDate := time.Now()
	order.BoughtDate = bgDate.String()
	order.TravelDate = trDate
	order.TrainNumber = info.TripId
	order.SeatClass = info.SeatType
	order.Status = uint16(util.Change)
	order.Price = ticketPrice

	if (rsi.IsTripGD(oldTripId) && rsi.IsTripGD(info.TripId)) || (!rsi.IsTripGD(oldTripId) && !rsi.IsTripGD(info.TripId)) {

		if rsi.IsTripGD(info.TripId) {
			order, err = rsi.orderService.Get().SaveOrderInfo(ctx, order, token)
		} else {
			order, err = rsi.orderOtherService.Get().SaveOrderInfo(ctx, order, token)
		}
		if err != nil {
			return util.Order{}, err
		}
	} else {

		if rsi.IsTripGD(oldTripId) {
			_, err = rsi.orderService.Get().DeleteOrder(ctx, order.Id, token)
		} else {
			_, err = rsi.orderOtherService.Get().DeleteOrder(ctx, order.Id, token)
		}

		if err != nil {
			return util.Order{}, err
		}

		if rsi.IsTripGD(info.TripId) {
			order, err = rsi.orderService.Get().CreateNewOrder(ctx, order, token)
		} else {
			order, err = rsi.orderOtherService.Get().CreateNewOrder(ctx, order, token)
		}
		if err != nil {
			return util.Order{}, err
		}
	}

	return order, nil
}
