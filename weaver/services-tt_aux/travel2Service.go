package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type Travel2Service interface {
	GetTrainTypeByTripId(ctx context.Context, tripId, token string) (util.Train, error)
	GetRouteByTripId(ctx context.Context, tripId, token string) (util.Route, error)
	GetTripsByRouteId(ctx context.Context, routeIds []string, token string) ([]util.Trip, error)
	UpdateTrip(ctx context.Context, trip util.Trip, token string) (util.Trip, error)
	Retrieve(ctx context.Context, tripId string, token string) (util.Trip, error)
	CreateTrip(ctx context.Context, trip util.Trip, token string) (util.Trip, error)
	DeleteTrip(ctx context.Context, tripId string, token string) (string, error)
	QueryInfo(ctx context.Context, startingPlace, endPlace, departureTime, token string) ([]util.TripDetails, error)
	GetTripAllDetailInfo(ctx context.Context, id, from, to, travelDate, token string) (util.Trip, util.TripDetails, error)
	GetTickets(ctx context.Context, id, from, to, travelDate, token string) (util.TripDetails, error)
	QueryAll(ctx context.Context, token string) ([]util.TripDetails, error)
	AdminQueryAll(ctx context.Context, token string) ([]util.Trip, []util.Train, []util.Route, error)
}

type travel2Service struct {
	weaver.Implements[Travel2Service]
	db                components.NoSQLDatabase
	routeService      weaver.Ref[RouteService]
	trainService      weaver.Ref[TrainService]
	ticketInfoService weaver.Ref[TicketInfoService]
	orderOtherService weaver.Ref[OrderOtherService]
	seatService       weaver.Ref[SeatService]
	roles             []string
}

func (tsi *travel2Service) GetTrainTypeByTripId(ctx context.Context, tripId, token string) (util.Train, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Train{}, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	//! TODO gotta verify this
	query := fmt.Sprintf(`{"Id": {"Type": %s, "Number": %s}}`, tripId[:1], tripId[1:])
	//! Can also do like this:
	// query := fmt.Sprintf(`{"Id": %v }`, TripId{
	// 	Type: tripId[:1]
	// 	Number: tripId[1:]
	// })

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Train{}, err
	}

	var trip util.Trip
	err = result.Decode(&trip)
	if err != nil {
		return util.Train{}, nil
	}

	trainType, err := tsi.trainService.Get().Retrieve(ctx, trip.TrainTypeId, token)
	if err != nil {
		return util.Train{}, err
	}

	return trainType, nil
}

func (tsi *travel2Service) GetRouteByTripId(ctx context.Context, tripId, token string) (util.Route, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Route{}, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	//! TODO gotta verify this; see previous route
	query := fmt.Sprintf(`{"Id": {"Type": %s, "Number": %s}}`, tripId[:1], tripId[1:])

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Route{}, err
	}

	var trip util.Trip

	err = result.Decode(&trip)
	if err != nil {
		return util.Route{}, err
	}

	route, err := tsi.routeService.Get().QueryById(ctx, trip.RouteId, token)
	if err != nil {
		return util.Route{}, err
	}
	return route, nil
}

func (tsi *travel2Service) GetTripsByRouteId(ctx context.Context, routeIds []string, token string) ([]util.Trip, error) {

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	var tripList []util.Trip
	for _, routeId := range routeIds {
		query := fmt.Sprintf(`{"RouteId": %s}`, routeId)

		result, err := collection.FindMany(query)
		if err != nil {
			fmt.Print(err)
			continue
		}

		var trips []util.Trip
		result.All(&trips)
		tripList = append(tripList, trips...)
	}

	return tripList, nil
}

func (tsi *travel2Service) UpdateTrip(ctx context.Context, trip util.Trip, token string) (util.Trip, error) {
	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return util.Trip{}, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	query := fmt.Sprintf(`{"Id": {"Type": %s, "Number": %s}}`, trip.TrainTypeId, trip.Number)

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Trip{}, err
	}

	var existingTrip util.Trip
	result.Decode(&existingTrip)
	if existingTrip.Id == "" {
		return util.Trip{}, errors.New("util.Trip not found.")
	}

	err = collection.ReplaceOne(query, trip)
	if err != nil {
		return util.Trip{}, err
	}

	return trip, nil
}

func (tsi *travel2Service) Retrieve(ctx context.Context, tripId string, token string) (util.Trip, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Trip{}, err
	}

	query := fmt.Sprintf(`{"Id": %s`, tripId)
	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Trip{}, err
	}

	var existingTrip util.Trip
	err = result.Decode(&existingTrip)
	if err != nil {
		return util.Trip{}, nil
	}

	return existingTrip, nil
}

func (tsi *travel2Service) CreateTrip(ctx context.Context, trip util.Trip, token string) (util.Trip, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Trip{}, err
	}

	query := fmt.Sprintf(`{"Id": {"Type": %s, "Number": %s}}`, trip.TrainTypeId, trip.Number)
	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Trip{}, err
	}

	var existingTrip util.Trip
	err = result.Decode(&existingTrip)
	if err != nil || existingTrip.Id == "" {
		return util.Trip{}, nil
	}

	err = collection.InsertOne(trip)
	if err != nil {
		return util.Trip{}, nil
	}

	return trip, nil
}

func (tsi *travel2Service) DeleteTrip(ctx context.Context, tripId string, token string) (string, error) {
	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf(`{"Id": %s}`, tripId)
	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Trip deleted.", nil
}

func (tsi *travel2Service) QueryInfo(ctx context.Context, startingPlace, endPlace, departureTime, token string) ([]util.TripDetails, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	startingPlaceId, err := tsi.ticketInfoService.Get().QueryForStationId(ctx, startingPlace, token)
	if err != nil {
		return nil, err
	}

	endPlaceId, err := tsi.ticketInfoService.Get().QueryForStationId(ctx, endPlace, token)
	if err != nil {
		return nil, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	result, err := collection.FindMany("")
	if err != nil {
		return nil, err
	}

	var trips []util.Trip
	result.All(&trips)

	var route util.Route

	var tripDetailsList []util.TripDetails
	var tripDetails util.TripDetails
	for _, trip := range trips {
		route, err = tsi.routeService.Get().QueryById(ctx, trip.RouteId, token)
		if err != nil {
			continue
		}

		foundLeft := false
		for _, station := range route.Stations {
			if startingPlaceId == station {
				foundLeft = true
				continue
			}

			if endPlaceId == station {
				if foundLeft {
					tripDetails, err = tsi.GetTickets(ctx, trip, route, startingPlaceId, endPlaceId, startingPlace, endPlace, departureTime, token)
					if err != nil {
						break
					}
					tripDetailsList = append(tripDetailsList, tripDetails)
				} else {
					break
				}
			}
		}
	}

	return tripDetailsList, nil

}

func (tsi *travel2Service) GetTripAllDetailInfo(ctx context.Context, id, from, to, travelDate, token string) (util.Trip, util.TripDetails, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Trip{}, util.TripDetails{}, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	//! TODO gotta verify this; see previous route
	query := fmt.Sprintf(`{"Id": {"Type": %s, "Number": %s}}`, id[:1], id[1:])

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Trip{}, util.TripDetails{}, err
	}

	var trip util.Trip
	result.Decode(&trip)

	startingPlaceId, err := tsi.ticketInfoService.Get().QueryForStationId(ctx, from, token)
	if err != nil {
		return util.Trip{}, util.TripDetails{}, err
	}

	endPlaceId, err := tsi.ticketInfoService.Get().QueryForStationId(ctx, to, token)
	if err != nil {
		return util.Trip{}, util.TripDetails{}, err
	}

	route, err := tsi.routeService.Get().QueryById(ctx, trip.RouteId, token)
	if err != nil {
		return util.Trip{}, util.TripDetails{}, err
	}

	tripResponse, err := tsi.GetTickets(ctx, trip, route, startingPlaceId, endPlaceId, from, to, travelDate, token)
	if err != nil {
		return util.Trip{}, util.TripDetails{}, err
	}

	return trip, tripResponse, nil
}

func (tsi *travel2Service) QueryAll(ctx context.Context, token string) ([]util.Trip, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	result, err := collection.FindMany("") //TODO verify
	if err != nil {
		return nil, err
	}

	var trips []util.Trip
	err = result.All(&trips)
	if err != nil {
		return nil, err
	}

	return trips, nil
}

func (tsi *travel2Service) AdminQueryAll(ctx context.Context, token string) ([]util.Trip, []util.Train, []util.Route, error) {
	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return nil, nil, nil, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trips")

	result, err := collection.FindMany("") //TODO verify
	if err != nil {
		return nil, nil, nil, err
	}

	var trips []util.Trip
	err = result.All(&trips)
	if err != nil {
		return nil, nil, nil, err
	}

	var routes []util.Route
	var trainTypes []util.Train

	for _, trip := range trips {

		route, _ := tsi.routeService.Get().QueryById(ctx, trip.RouteId, token)
		tt, _ := tsi.trainService.Get().Retrieve(ctx, trip.TrainTypeId, token)

		routes = append(routes, route)
		trainTypes = append(trainTypes, tt)
	}

	return trips, trainTypes, routes, nil
}

//*************************************************************************************

func (tsi *travel2Service) GetTickets(ctx context.Context, trip util.Trip, route util.Route, startingPlaceId, endPlaceId, startingPlaceName, endPlaceName, departureTime, token string) (util.TripDetails, error) {

	dateFormat := "Sat Jul 26 00:00:00 2025"

	depart, _ := time.Parse(dateFormat, departureTime)

	if depart.After(time.Now()) {
		return util.TripDetails{}, errors.New("Departure time in the past.")
	}

	resForTravel, err := tsi.ticketInfoService.Get().QueryForTravel(ctx, util.Travel{
		Trip:          trip,
		StartingPlace: startingPlaceName,
		EndPlace:      endPlaceName,
		DepartureTime: departureTime,
	}, token)

	if err != nil {
		return util.TripDetails{}, err
	}

	soldTicket, err := tsi.orderOtherService.Get().CalculateSoldTicket(ctx, departureTime, trip.TrainTypeId+trip.Number, token)
	if err != nil {
		return util.TripDetails{}, err
	}

	seat := util.Seat{
		TravelDate:   departureTime,
		TrainNumber:  trip.TrainTypeId + trip.Number,
		StartStation: startingPlaceId,
		DestStation:  endPlaceId,
		SeatType:     uint16(util.FirstClass),
	}

	//! TODO revisit after seatService impl
	first, err := tsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
	if err != nil {
		return util.TripDetails{}, err
	}
	seat.SeatType = uint16(util.SecondClass)

	//! TODO revisit after seatService impl
	second, err := tsi.seatService.Get().GetLeftTicketOfInterval(ctx, seat, token)
	if err != nil {
		return util.TripDetails{}, err
	}

	trainType, err := tsi.trainService.Get().Retrieve(ctx, trip.TrainTypeId, token)
	if err != nil {
		return util.TripDetails{}, err
	}

	tripResponse := util.TripDetails{
		ComfortClass:    first,
		EconomyClass:    second,
		StartingStation: startingPlaceName,
		EndStation:      endPlaceName,
	}

	var indexStart int
	var indexEnd int
	for idx, st := range route.Stations {

		if st == startingPlaceId {
			indexStart = idx
		}
		if st == endPlaceId {
			indexEnd = idx
		}
	}

	distanceStart := route.Distances[indexStart] - route.Distances[0]
	distanceEnd := route.Distances[indexEnd] - route.Distances[0]

	minutesStart := 60 * distanceStart / trainType.AvgSpeed
	minutesEnd := 60 * distanceEnd / trainType.AvgSpeed

	tmpTime, _ := time.Parse(dateFormat, trip.StartingTime)
	startingTime := tmpTime.Add(time.Minute * time.Duration(minutesStart))
	endTime := tmpTime.Add(time.Minute * time.Duration(minutesEnd))

	tripResponse.StartingTime = startingTime.String()
	tripResponse.EndTime = endTime.String()

	tripResponse.TripId = soldTicket.TrainNumber
	tripResponse.TrainTypeId = trip.TrainTypeId
	tripResponse.PriceForComfortClass = resForTravel.Prices["ComfortClass"]
	tripResponse.PriceForEconomyClass = resForTravel.Prices["EconomyClass"]

	return tripResponse, nil
}
