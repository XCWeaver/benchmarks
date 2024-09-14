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

//! *****************************************************************************************************************************************
//! **                                                                                                                                     **
//! **          This service uses a queue to push email notifications, but we found it to be disabled in the original implementation 	   **
//! **          and thus removed it.   																									   **
//! **                                                                                                                                     **
//! *****************************************************************************************************************************************

type PreserveOtherService interface {
	Preserve(ctx context.Context, info util.OrderInfo, token string) (string, error)
}

type preserveOtherService struct {
	weaver.Implements[PreserveOtherService]
	securityService   weaver.Ref[SecurityService]
	contactService    weaver.Ref[ContactService]
	travel2Service    weaver.Ref[Travel2Service]
	stationService    weaver.Ref[StationService]
	foodService       weaver.Ref[FoodService]
	insuranceService  weaver.Ref[InsuranceService]
	consignService    weaver.Ref[ConsignService]
	orderOtherService weaver.Ref[OrderOtherService]
	seatService       weaver.Ref[SeatService]
	userService       weaver.Ref[UserService]
	ticketInfoService weaver.Ref[TicketInfoService]
	roles             []string
}

func (psi *preserveOtherService) Preserve(ctx context.Context, info util.OrderTicketInfo, token string) (string, error) {
	err := util.Authenticate(token, psi.roles...)
	if err != nil {
		return "", err
	}

	var err1, err2, err3 error
	var contact util.Contact
	var trip util.Trip
	var tripData util.TripDetails

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		_, err1 = psi.securityService.Get().Check(ctx, info.AccountId, token)
	}()

	go func() {
		defer wg.Done()
		contact, err2 = psi.contactService.Get().GetContactsByContactId(ctx, info.TripId, token)
	}()

	go func() {
		defer wg.Done()
		trip, tripData, err3 = psi.travel2Service.Get().GetTripAllDetailInfo(ctx, info.TripId, info.From, info.To, info.Date, token)
	}()

	wg.Wait()

	if err1 != nil {
		return "", err1
	}
	if err2 != nil {
		return "", err2
	}
	if err3 != nil {
		return "", err3
	}

	//* ****************************************************************************************
	if info.SeatType == uint16(util.FirstClass) && tripData.ComfortClass == 0 {
		return "", errors.New("Not enough seats to perform this booking.")
	} else if tripData.EconomyClass == 0 {
		return "", errors.New("Not enough seats to perform this booking.")
	}

	wg.Add(3)

	var fromStationId, toStationId string
	var travelResult util.TravelResult
	go func() {
		defer wg.Done()
		fromStationId, err1 = psi.stationService.Get().QueryForStationId(ctx, info.From, token)
	}()

	go func() {
		defer wg.Done()
		toStationId, err2 = psi.stationService.Get().QueryForStationId(ctx, info.To, token)
	}()

	go func() {
		defer wg.Done()
		travel := util.Travel{
			Trip:          trip,
			StartingPlace: info.From,
			EndPlace:      info.To,
			DepartureTime: time.Now().Format(time.ANSIC),
		}
		travelResult, err3 = psi.ticketInfoService.Get().QueryForTravel(ctx, travel, token)
	}()
	wg.Wait()

	if err1 != nil {
		return "", err1
	}
	if err2 != nil {
		return "", err2
	}
	if err3 != nil {
		return "", err3
	}

	//* ****************************************************************************************

	var ticket util.Ticket
	var orderSeatClass uint16
	var orderPrice float32
	if info.SeatType == uint16(util.FirstClass) {
		orderSeatClass = 1
		orderPrice = travelResult.Prices["ComfortClass"]
		ticket, err = psi.seatService.Get().Create(ctx, 1, info.Date, info.TripId, fromStationId, toStationId, token)
	} else {
		orderSeatClass = 2
		orderPrice = travelResult.Prices["EconomyClass"]
		ticket, err = psi.seatService.Get().Create(ctx, 2, info.Date, info.TripId, fromStationId, toStationId, token)
	}
	if err != nil {
		return "", err
	}

	orderSeatNum := ticket.SeatNo

	tmpDate, _ := time.Parse(time.ANSIC, info.Date)
	y, m, d := tmpDate.Date()
	startTime, _ := time.Parse(time.ANSIC, tripData.StartingTime)
	hr, min, sec := startTime.Clock()

	order := util.Order{
		BoughtDate:             time.Now().Format(time.ANSIC),
		TravelDate:             fmt.Sprintf("%s %s %s %s:%s:%s %s", tmpDate.Weekday().String()[:3], m.String()[:3], d, hr, min, sec, y),
		AccountId:              info.AccountId,
		ContactsName:           info.ContactsId,
		DocumentType:           contact.DocumentType,
		ContactsDocumentNumber: contact.DocumentNumber,
		TrainNumber:            info.TripId,
		SeatClass:              orderSeatClass,
		SeatNumber:             orderSeatNum,
		From:                   fromStationId,
		To:                     toStationId,
		Status:                 0,
		Price:                  orderPrice,
	}

	order, err = psi.orderOtherService.Get().CreateNewOrder(ctx, order, token)
	if err != nil {
		return "", err
	}

	if info.Insurance != 0 {
		_, err = psi.insuranceService.Get().CreateNewInsurance(ctx, info.Insurance, order.Id, token)
	}
	if err != nil {
		return "", err
	}

	if info.FoodType != 0 {
		foodOrder := util.FoodOrder{
			OrderId:     order.Id,
			FoodType:    info.FoodType,
			StationName: info.StationName,
			StoreName:   info.StoreName,
			FoodName:    info.FoodName,
			Price:       info.FoodPrice,
		}

		_, err = psi.foodService.Get().CreateFoodOrder(ctx, foodOrder, token)
	}
	if err != nil {
		return "", err
	}

	if info.ConsigneeName != "" {
		consign := util.Consign{
			OrderId:    order.Id,
			AccountId:  order.AccountId,
			HandleDate: info.HandleDate,
			TargetDate: order.TravelDate,
			From:       order.From,
			To:         order.To,
			Consignee:  info.ConsigneeName,
			Phone:      info.ConsigneePhone,
			Weight:     info.ConsigneWeight,
			Within:     info.IsWithin,
		}

		_, err = psi.consignService.Get().InsertConsign(ctx, consign, token)
	}
	if err != nil {
		return "", err
	}

	//* Fetching the user would be needed if we use the email-queue logic
	//* We still make the function call for sameness' sake
	_, err = psi.userService.Get().GetUserById(ctx, info.AccountId, token)
	if err != nil {
		return "", err
	}

	return "Booking successful", nil
}
