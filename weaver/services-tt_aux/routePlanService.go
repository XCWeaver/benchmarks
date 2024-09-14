package services

import (
	"context"
	"math"
	"sync"
	"time"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type RoutePlanService interface {
	GetCheapestRoutes(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error)
	GetQuickestRoutes(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error)
	GetMinStopStations(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error)
}

type routePlanService struct {
	weaver.Implements[RoutePlanService]
	travelService     weaver.Ref[TravelService]
	travel2Service    weaver.Ref[Travel2Service]
	seatService       weaver.Ref[SeatService]
	stationService    weaver.Ref[StationService]
	ticketInfoService weaver.Ref[TicketInfoService]
	routeService      weaver.Ref[RouteService]
	// roles []string
}

func (rpsi *routePlanService) GetCheapestRoutes(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	var tripsFirst []util.TripDetails
	var tripsSecond []util.TripDetails
	var modData, returnData []util.TripDetails

	go func() {
		defer wg.Done()
		tripsFirst, err1 = rpsi.travelService.Get().QueryInfo(ctx, info.FromStationName, info.ToStationName, info.TravelDate, token)
	}()

	go func() {
		defer wg.Done()
		tripsSecond, err2 = rpsi.travel2Service.Get().QueryInfo(ctx, info.FromStationName, info.ToStationName, info.TravelDate, token)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	allTrips := append(tripsFirst, tripsSecond...)

	_range := 5
	minPrice := float32(0.0)
	minIndex := -1

	if len(allTrips) < _range {
		_range = len(allTrips)
	}

	//* find <5 cheapest
	for i := 0; i < _range; i++ {
		minPrice = math.MaxFloat32

		for j := 0; j < len(allTrips); j++ {
			tripResp := allTrips[j]

			if tripResp.PriceForEconomyClass < minPrice {
				minPrice = tripResp.PriceForEconomyClass
				minIndex = j
			}
		}

		modData = append(modData, allTrips[minIndex])
		//remove this item
		//preserving order: allTrips = append(allTrips[:minIndex], allTrips[minIndex+1:]...)
		allTrips[minIndex] = allTrips[len(allTrips)-1]
		allTrips = allTrips[:len(allTrips)-1]
	}

	var route util.Route
	for _, selectedTrip := range modData {
		var err error
		if selectedTrip.TripId[:1] == "G" || selectedTrip.TripId[:1] == "D" {
			route, err = rpsi.travelService.Get().GetRouteByTripId(ctx, selectedTrip.TripId, token)
		} else {
			route, err = rpsi.travel2Service.Get().GetRouteByTripId(ctx, selectedTrip.TripId, token)
		}
		if err != nil {
			return returnData, err
		}

		selectedTrip.StopStations = route.Stations
		returnData = append(returnData, selectedTrip)
	}

	return returnData, nil
}

func (rpsi *routePlanService) GetQuickestRoutes(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var err1, err2 error
	var tripsFirst []util.TripDetails
	var tripsSecond []util.TripDetails
	var modData, returnData []util.TripDetails

	go func() {
		defer wg.Done()
		tripsFirst, err1 = rpsi.travelService.Get().QueryInfo(ctx, info.FromStationName, info.ToStationName, info.TravelDate, token)
	}()

	go func() {
		defer wg.Done()
		tripsSecond, err2 = rpsi.travel2Service.Get().QueryInfo(ctx, info.FromStationName, info.ToStationName, info.TravelDate, token)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	allTrips := append(tripsFirst, tripsSecond...)

	_range := 5
	var minTime int64
	minIndex := -1

	if len(allTrips) < _range {
		_range = len(allTrips)
	}

	//* find <5 shortest
	for i := 0; i < _range; i++ {
		minTime = math.MaxInt64

		for j := 0; j < len(allTrips); j++ {
			tripResp := allTrips[j]

			//TODO rework
			endTime, _ := time.Parse(time.ANSIC, tripResp.EndTime)
			startTime, _ := time.Parse(time.ANSIC, tripResp.StartingTime)
			timeDiff := endTime.Sub(startTime).Milliseconds()

			if timeDiff < minTime {
				minTime = timeDiff
				minIndex = j
			}
		}

		modData = append(modData, allTrips[minIndex])
		//remove this item
		//preserving order: allTrips = append(allTrips[:minIndex], allTrips[minIndex+1:]...)
		allTrips[minIndex] = allTrips[len(allTrips)-1]
		allTrips = allTrips[:len(allTrips)-1]
	}

	var route util.Route
	for _, selectedTrip := range modData {
		var err error
		if selectedTrip.TripId[:1] == "G" || selectedTrip.TripId[:1] == "D" {
			route, err = rpsi.travelService.Get().GetRouteByTripId(ctx, selectedTrip.TripId, token)
		} else {
			route, err = rpsi.travel2Service.Get().GetRouteByTripId(ctx, selectedTrip.TripId, token)
		}
		if err != nil {
			return returnData, err
		}
		selectedTrip.StopStations = route.Stations
		returnData = append(returnData, selectedTrip)
	}

	return returnData, nil
}

func (rpsi *routePlanService) GetMinStopStations(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(3)
	var err1, err2, err3 error
	var fromStationId, toStationId string
	var routes []util.Route

	go func() {
		defer wg.Done()
		fromStationId, err1 = rpsi.stationService.Get().QueryForStationId(ctx, info.FromStationName, token)
	}()

	go func() {
		defer wg.Done()
		toStationId, err2 = rpsi.stationService.Get().QueryForStationId(ctx, info.ToStationName, token)
	}()

	go func() {
		defer wg.Done()
		routes, err3 = rpsi.routeService.Get().QueryByStartAndTerminal(ctx, info.FromStationName, info.ToStationName, token)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	if err3 != nil {
		return nil, err3
	}

	var gapList []int

	for _, route := range routes {
		indexStart := 0
		indexEnd := 0

		for idx, s := range route.Stations {
			if s == fromStationId {
				indexStart = idx
				continue
			}
			if s == toStationId {
				indexEnd = idx
			}
		}

		gapList = append(gapList, indexEnd-indexStart)
	}

	_range := 5
	minIndex := -1

	var resultRoutes []string
	if len(routes) < _range {
		_range = len(routes)
	}

	for i := 0; i < _range; i++ {

		minGap := math.MaxInt

		for j := 0; j < len(gapList); j++ {
			if gapList[j] < minGap {
				minGap = gapList[j]
				minIndex = j
			}
		}

		resultRoutes = append(resultRoutes, routes[minIndex].Id)

		routes[minIndex] = routes[len(routes)-1]
		routes = routes[:len(routes)-1]

		gapList[minIndex] = gapList[len(gapList)-1]
		gapList = gapList[:len(gapList)-1]
	}

	//******************** Now get the actual trips w/ details *************************

	var wg2 sync.WaitGroup
	wg2.Add(2)
	var firstTrips, secondTrips []util.Trip
	go func() {
		defer wg2.Done()
		firstTrips, err1 = rpsi.travelService.Get().GetTripsByRouteId(ctx, resultRoutes, token)
	}()
	go func() {
		defer wg2.Done()
		secondTrips, err2 = rpsi.travel2Service.Get().GetTripsByRouteId(ctx, resultRoutes, token)
	}()
	wg2.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	allTrips := append(firstTrips, secondTrips...)

	var finalTripsDetailed []util.TripDetails
	var tripDetails util.TripDetails

	for _, trip := range allTrips {

		if trip.Id[:1] == "G" || trip.Id[:1] == "D" {
			_, tripDetails, err = rpsi.travelService.Get().GetTripAllDetailInfo(ctx, trip.Id, info.FromStationName, info.ToStationName, info.TravelDate, token)
		} else {
			_, tripDetails, err = rpsi.travel2Service.Get().GetTripAllDetailInfo(ctx, trip.Id, info.FromStationName, info.ToStationName, info.TravelDate, token)
		}

		if err != nil {
			continue
		}

		route, err := rpsi.routeService.Get().QueryById(ctx, trip.RouteId, token)
		if err != nil {
			continue
		}
		tripDetails.StopStations = route.Stations

		finalTripsDetailed = append(finalTripsDetailed, tripDetails)
	}

	return finalTripsDetailed, nil
}
