package services

import (
	"context"
	"sync"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type TravelPlanService interface {
	GetTransferResult(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, []util.TripDetails, error)
	GetByCheapest(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error)
	GetByQuickest(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error)
	GetByMinStation(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error)
}

type travelPlanService struct {
	weaver.Implements[TravelPlanService]
	travelService     weaver.Ref[TravelService]
	travel2Service    weaver.Ref[Travel2Service]
	seatService       weaver.Ref[SeatService]
	stationService    weaver.Ref[StationService]
	ticketInfoService weaver.Ref[TicketInfoService]
	routePlanService  weaver.Ref[RoutePlanService]
	// roles []string
}

func (tpsi *travelPlanService) GetTransferResult(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, []util.TripDetails, error) {

	err := util.Authenticate(token)
	if err != nil {
		return nil, nil, err
	}

	var wg sync.WaitGroup
	wg.Add(4)

	var err1, err2, err3, err4 error
	var firstSectionFromHighSpeed, firstSectionFromNormal, secondSectionFromHighSpeed, secondSectionFromNormal []util.TripDetails

	go func() {
		defer wg.Done()
		firstSectionFromHighSpeed, err1 = tpsi.travelService.Get().QueryInfo(ctx, info.FromStationName, info.ViaStationName, info.TravelDate, token)
	}()

	go func() {
		defer wg.Done()
		firstSectionFromNormal, err2 = tpsi.travel2Service.Get().QueryInfo(ctx, info.FromStationName, info.ViaStationName, info.TravelDate, token)
	}()

	go func() {
		defer wg.Done()
		secondSectionFromHighSpeed, err1 = tpsi.travelService.Get().QueryInfo(ctx, info.ViaStationName, info.ToStationName, info.TravelDate, token)
	}()

	go func() {
		defer wg.Done()
		secondSectionFromNormal, err2 = tpsi.travel2Service.Get().QueryInfo(ctx, info.ViaStationName, info.ToStationName, info.TravelDate, token)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, nil, err1
	}
	if err2 != nil {
		return nil, nil, err2
	}
	if err3 != nil {
		return nil, nil, err3
	}
	if err4 != nil {
		return nil, nil, err4
	}

	firstSection := append(firstSectionFromHighSpeed, firstSectionFromNormal...)
	secondSection := append(secondSectionFromHighSpeed, secondSectionFromNormal...)

	return firstSection, secondSection, nil
}

func (tpsi *travelPlanService) GetByCheapest(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	cheapestTrips, err := tpsi.routePlanService.Get().GetCheapestRoutes(ctx, info, token)
	if err != nil {
		return nil, err
	}

	var finalTrips []util.TripDetails

	var err1, err2, err3, err4, err5 error
	var wg sync.WaitGroup
	var stopNames []string
	var fromStationId, toStationId string
	var first, second uint16
	for _, trip := range cheapestTrips {

		wg.Add(3)
		go func() {
			defer wg.Done()
			stopNames, err1 = tpsi.stationService.Get().QueryForNameBatch(ctx, trip.StopStations, token)
		}()
		go func() {
			defer wg.Done()
			fromStationId, err2 = tpsi.stationService.Get().QueryForStationId(ctx, trip.StartingStation, token)
		}()
		go func() {
			defer wg.Done()
			toStationId, err3 = tpsi.stationService.Get().QueryForStationId(ctx, trip.EndStation, token)
		}()
		wg.Wait()

		seat := util.Seat{
			TravelDate:   info.TravelDate,
			TrainNumber:  trip.TripId,
			StartStation: fromStationId,
			DestStation:  toStationId,
			SeatType:     uint16(util.FirstClass),
		}

		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}

		wg.Add(2)
		go func() {
			defer wg.Done()
			first, err4 = tpsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
		}()

		go func() {
			seat.SeatType = uint16(util.SecondClass)
			defer wg.Done()
			second, err5 = tpsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
		}()
		wg.Wait()

		if err4 != nil || err5 != nil {
			continue
		}

		trip.StopStations = stopNames
		trip.NumberOfRestTicketFirstClass = first
		trip.NumberOfRestTicketSecondClass = second

		finalTrips = append(finalTrips, trip)
	}

	return finalTrips, nil
}

func (tpsi *travelPlanService) GetByQuickest(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	quickestTrips, err := tpsi.routePlanService.Get().GetQuickestRoutes(ctx, info, token)
	if err != nil {
		return nil, err
	}

	var finalTrips []util.TripDetails

	var err1, err2, err3, err4, err5 error
	var wg sync.WaitGroup
	var stopNames []string
	var fromStationId, toStationId string
	var first, second uint16
	for _, trip := range quickestTrips {

		wg.Add(3)
		go func() {
			defer wg.Done()
			stopNames, err1 = tpsi.stationService.Get().QueryForNameBatch(ctx, trip.StopStations, token)
		}()
		go func() {
			defer wg.Done()
			fromStationId, err2 = tpsi.stationService.Get().QueryForStationId(ctx, trip.StartingStation, token)
		}()
		go func() {
			defer wg.Done()
			toStationId, err3 = tpsi.stationService.Get().QueryForStationId(ctx, trip.EndStation, token)
		}()
		wg.Wait()

		seat := util.Seat{
			TravelDate:   info.TravelDate,
			TrainNumber:  trip.TripId,
			StartStation: fromStationId,
			DestStation:  toStationId,
			SeatType:     uint16(util.FirstClass),
		}

		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}

		wg.Add(2)
		go func() {
			defer wg.Done()
			first, err4 = tpsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
		}()

		go func() {
			seat.SeatType = uint16(util.SecondClass)
			defer wg.Done()
			second, err5 = tpsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
		}()
		wg.Wait()

		if err4 != nil || err5 != nil {
			continue
		}

		trip.StopStations = stopNames
		trip.NumberOfRestTicketFirstClass = first
		trip.NumberOfRestTicketSecondClass = second

		finalTrips = append(finalTrips, trip)
	}

	return finalTrips, nil

}

func (tpsi *travelPlanService) GetByMinStation(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	lessHopTripsquickestTrips, err := tpsi.routePlanService.Get().GetMinStopStations(ctx, info, token)
	if err != nil {
		return nil, err
	}

	var finalTrips []util.TripDetails

	var err1, err2, err3, err4, err5 error
	var wg sync.WaitGroup
	var stopNames []string
	var fromStationId, toStationId string
	var first, second uint16
	for _, trip := range lessHopTripsquickestTrips {

		wg.Add(3)
		go func() {
			defer wg.Done()
			stopNames, err1 = tpsi.stationService.Get().QueryForNameBatch(ctx, trip.StopStations, token)
		}()
		go func() {
			defer wg.Done()
			fromStationId, err2 = tpsi.stationService.Get().QueryForStationId(ctx, trip.StartingStation, token)
		}()
		go func() {
			defer wg.Done()
			toStationId, err3 = tpsi.stationService.Get().QueryForStationId(ctx, trip.EndStation, token)
		}()
		wg.Wait()

		seat := util.Seat{
			TravelDate:   info.TravelDate,
			TrainNumber:  trip.TripId,
			StartStation: fromStationId,
			DestStation:  toStationId,
			SeatType:     uint16(util.FirstClass),
		}

		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}

		wg.Add(2)
		go func() {
			defer wg.Done()
			first, err4 = tpsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
		}()

		go func() {
			seat.SeatType = uint16(util.SecondClass)
			defer wg.Done()
			second, err5 = tpsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
		}()
		wg.Wait()

		if err4 != nil || err5 != nil {
			continue
		}

		trip.StopStations = stopNames
		trip.NumberOfRestTicketFirstClass = first
		trip.NumberOfRestTicketSecondClass = second

		finalTrips = append(finalTrips, trip)
	}

	return finalTrips, nil

}
