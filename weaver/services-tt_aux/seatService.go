package services

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type SeatService interface {
	Create(ctx context.Context, seatType uint16, travelDate, trainNumber, startStation, destStation, token string) (util.Ticket, error)
	GetLeftTicketOfInterval(ctx context.Context, seat util.Seat, token string) (uint16, error)
}

type seatService struct {
	weaver.Implements[SeatService]
	travelService     weaver.Ref[TravelService]
	travel2Service    weaver.Ref[Travel2Service]
	orderService      weaver.Ref[OrderService]
	orderOtherService weaver.Ref[OrderOtherService]
	configService     weaver.Ref[ConfigService]
	roles             []string
}

func (ssi *seatService) Create(ctx context.Context, seatType uint16, travelDate, trainNumber, startStation, destStation, token string) (util.Ticket, error) {
	err := util.Authenticate(token, ssi.roles[0])

	if err != nil {
		return util.Ticket{}, err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	var err1, err2, err3 error
	var route util.Route
	var leftTicketInfo []util.Ticket
	var trainType util.Train

	if trainNumber[:1] == "G" || trainNumber[:1] == "D" {
		go func() {
			defer wg.Done()
			route, err1 = ssi.travelService.Get().GetRouteByTripId(ctx, trainNumber, token)
		}()

		go func() {
			defer wg.Done()
			leftTicketInfo, err2 = ssi.orderService.Get().GetTicketListByDateAndTripId(ctx, travelDate, trainNumber, token)
		}()

		go func() {
			defer wg.Done()
			trainType, err3 = ssi.travelService.Get().GetTrainTypeByTripId(ctx, trainNumber, token)
		}()
		wg.Wait()
	} else {

		go func() {
			defer wg.Done()
			route, err1 = ssi.travel2Service.Get().GetRouteByTripId(ctx, trainNumber, token)
		}()

		go func() {
			defer wg.Done()
			leftTicketInfo, err2 = ssi.orderOtherService.Get().GetTicketListByDateAndTripId(ctx, travelDate, trainNumber, token)
		}()

		go func() {
			defer wg.Done()
			trainType, err3 = ssi.travel2Service.Get().GetTrainTypeByTripId(ctx, trainNumber, token)
		}()
		wg.Wait()
	}

	if err1 != nil {
		return util.Ticket{}, err1
	}
	if err2 != nil {
		return util.Ticket{}, err2
	}
	if err3 != nil {
		return util.Ticket{}, err3
	}

	var seatTotalNum uint16

	if seatType == uint16(util.FirstClass) {
		seatTotalNum = trainType.ComfortClass
	} else {
		seatTotalNum = trainType.EconomyClass
	}

	finalTicket := util.Ticket{
		StartStation: startStation,
		DestStation:  destStation,
	}

	stationList := route.Stations

	var startStationIndex, destStationIndex int
	for idx, x := range stationList {
		if x == startStation {
			startStationIndex = idx
		}
	}
	for _, ticket := range leftTicketInfo {

		destStationIndex = 0
		for idx, x := range stationList {
			if x == ticket.DestStation {
				destStationIndex = idx
				break
			}
		}

		if destStationIndex < startStationIndex {
			finalTicket.SeatNo = ticket.SeatNo
			return finalTicket, nil
		}
	}

	seed := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(seed)

	genSeatNo := randomizer.Intn(int(seatTotalNum-1)) + 1

	for {
		if finalTicket.SeatNo == fmt.Sprintf("%d", genSeatNo) {
			genSeatNo = randomizer.Intn(int(seatTotalNum-1)) + 1
			continue
		}
		break
	}

	finalTicket.SeatNo = fmt.Sprintf("%d", genSeatNo)

	return finalTicket, nil
}

func (ssi *seatService) GetLeftTicketOfInterval(ctx context.Context, seat util.Seat, token string) (uint16, error) {
	err := util.Authenticate(token, ssi.roles[0])

	if err != nil {
		return 0, err
	}

	var wg sync.WaitGroup
	wg.Add(3)

	var err1, err2, err3 error
	var route util.Route
	var leftTicketInfo []util.Ticket
	var trainType util.Train

	if seat.TrainNumber[:1] == "G" || seat.TrainNumber[:1] == "D" {
		go func() {
			defer wg.Done()
			route, err1 = ssi.travelService.Get().GetRouteByTripId(ctx, seat.TrainNumber, token)
		}()

		go func() {
			defer wg.Done()
			leftTicketInfo, err2 = ssi.orderService.Get().GetTicketListByDateAndTripId(ctx, seat.TravelDate, seat.TrainNumber, token)
		}()

		go func() {
			defer wg.Done()
			trainType, err3 = ssi.travelService.Get().GetTrainTypeByTripId(ctx, seat.TrainNumber, token)
		}()
		wg.Wait()
	} else {

		go func() {
			defer wg.Done()
			route, err1 = ssi.travel2Service.Get().GetRouteByTripId(ctx, seat.TrainNumber, token)
		}()

		go func() {
			defer wg.Done()
			leftTicketInfo, err2 = ssi.orderOtherService.Get().GetTicketListByDateAndTripId(ctx, seat.TravelDate, seat.TrainNumber, token)
		}()

		go func() {
			defer wg.Done()
			trainType, err3 = ssi.travel2Service.Get().GetTrainTypeByTripId(ctx, seat.TrainNumber, token)
		}()
		wg.Wait()
	}

	if err1 != nil {
		return 0, err1
	}
	if err2 != nil {
		return 0, err2
	}
	if err3 != nil {
		return 0, err3
	}

	var seatTotalNum uint16

	if seat.SeatType == uint16(util.FirstClass) {
		seatTotalNum = trainType.ComfortClass
	} else {
		seatTotalNum = trainType.EconomyClass
	}

	stationList := route.Stations

	soldTicketCount := len(leftTicketInfo)

	var numTicketsLeft uint16

	var startStationIndex, destStationIndex int
	for idx, x := range stationList {
		if x == seat.StartStation {
			startStationIndex = idx
		}
	}
	for _, ticket := range leftTicketInfo {

		destStationIndex = 0
		for idx, x := range stationList {
			if x == ticket.DestStation {
				destStationIndex = idx
				break
			}
		}

		if destStationIndex < startStationIndex {
			numTicketsLeft += 1
		}
	}

	config, err := ssi.configService.Get().Retrieve(ctx, "DirectTicketAllocationProportion")

	if err != nil {
		return 0, nil
	}

	direstPart, _ := strconv.ParseFloat(config.Value, 32)

	if route.Stations[0] != seat.StartStation || route.Stations[len(route.Stations)-1] != seat.DestStation {
		direstPart = 1.0 - direstPart
	}

	unusedNum := uint16(float64(seatTotalNum)*direstPart - float64(soldTicketCount))

	numTicketsLeft += unusedNum

	return numTicketsLeft, nil
}
