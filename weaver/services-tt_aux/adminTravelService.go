package services

import (
	"context"
	"errors"
	"sync"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type AdminTravelService interface {
	GetAllTravels(ctx context.Context, token string) ([]util.Trip, []util.Train, []util.Route, error)
	AddTravel(ctx context.Context, trip util.Trip, token string) (util.Trip, error)
	UpdateTravel(ctx context.Context, trip util.Trip, token string) (util.Trip, error)
	DeleteTravel(ctx context.Context, tripId, token string) (string, error)
}

type adminTravelService struct {
	weaver.Implements[AdminTravelService]
	travelService  weaver.Ref[TravelService]
	travel2Service weaver.Ref[Travel2Service]
	roles          []string
}

func (atsi *adminTravelService) GetAllTravels(ctx context.Context, token string) ([]util.Trip, []util.Train, []util.Route, error) {

	err := util.Authenticate(token, atsi.roles...)
	if err != nil {
		return nil, nil, nil, err
	}

	var err1, err2 error
	var trip1, trip2 []util.Trip
	var trainType1, trainType2 []util.Train
	var route1, route2 []util.Route
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		trip1, trainType1, route1, err1 = atsi.travelService.Get().AdminQueryAll(ctx, token)
	}()

	go func() {
		defer wg.Done()
		trip2, trainType2, route2, err2 = atsi.travel2Service.Get().AdminQueryAll(ctx, token)
	}()

	wg.Wait()

	if err1 != nil {
		return nil, nil, nil, err1
	}
	if err2 != nil {
		return nil, nil, nil, err2
	}

	trips := append(trip1, trip2...)
	trainTypes := append(trainType1, trainType2...)
	routes := append(route1, route2...)

	if len(trips) == 0 {
		return nil, nil, nil, errors.New("No trips found")
	}

	return trips, trainTypes, routes, nil
}

func (atsi *adminTravelService) AddTravel(ctx context.Context, trip util.Trip, token string) (util.Trip, error) {

	err := util.Authenticate(token, atsi.roles...)
	if err != nil {
		return util.Trip{}, err
	}

	if trip.Id[0:1] == "D" || trip.Id[0:1] == "G" {
		trip, err = atsi.travelService.Get().CreateTrip(ctx, trip, token)
	} else {
		trip, err = atsi.travel2Service.Get().CreateTrip(ctx, trip, token)
	}

	if err != nil {
		return util.Trip{}, err
	}

	return trip, nil
}

func (atsi *adminTravelService) UpdateTravel(ctx context.Context, trip util.Trip, token string) (util.Trip, error) {

	err := util.Authenticate(token, atsi.roles...)
	if err != nil {
		return util.Trip{}, err
	}

	if trip.Id[0:1] == "D" || trip.Id[0:1] == "G" {
		trip, err = atsi.travelService.Get().UpdateTrip(ctx, trip, token)
	} else {
		trip, err = atsi.travel2Service.Get().UpdateTrip(ctx, trip, token)
	}

	if err != nil {
		return util.Trip{}, err
	}

	return trip, nil
}

func (atsi *adminTravelService) DeleteRoute(ctx context.Context, tripId, token string) (string, error) {
	err := util.Authenticate(token, atsi.roles...)
	if err != nil {
		return "", err
	}

	if tripId[0:1] == "D" || tripId[0:1] == "G" {
		_, err = atsi.travelService.Get().DeleteTrip(ctx, tripId, token)
	} else {
		_, err = atsi.travel2Service.Get().DeleteTrip(ctx, tripId, token)
	}
	if err != nil {
		return "", err
	}

	return "util.Trip deleted.", nil
}
