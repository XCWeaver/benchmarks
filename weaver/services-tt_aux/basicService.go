package services

import (
	"context"
	"sync"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type BasicService interface {
	QueryForTravel(ctx context.Context, info util.Travel, token string) (util.TravelResult, error)
	QueryForStationId(ctx context.Context, name, token string) (string, error)
}

type basicService struct {
	weaver.Implements[BasicService]
	trainService   weaver.Ref[TrainService]
	stationService weaver.Ref[StationService]
	routeService   weaver.Ref[RouteService]
	priceService   weaver.Ref[PriceService]
}

func (bsi *basicService) QueryForTravel(ctx context.Context, info util.Travel, token string) (util.TravelResult, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.TravelResult{}, err
	}

	var wg sync.WaitGroup
	wg.Add(5)

	var train util.Train
	var err1 error
	go func() {
		defer wg.Done()
		train, err1 = bsi.trainService.Get().Retrieve(ctx, info.Trip.TrainTypeId, token)
	}()

	var route util.Route
	var err2 error
	go func() {
		defer wg.Done()
		route, err2 = bsi.routeService.Get().QueryById(ctx, info.Trip.RouteId, token)
	}()

	var priceConfig util.PriceConfig
	var err3 error
	go func() {
		defer wg.Done()
		priceConfig, err3 = bsi.priceService.Get().Query(ctx, info.Trip.RouteId, info.Trip.TrainTypeId, token)
	}()

	var startingPlaceId string
	var err4 error
	go func() {
		defer wg.Done()
		startingPlaceId, err4 = bsi.stationService.Get().QueryForStationId(ctx, info.StartingPlace, token)
	}()

	var endPlaceId string
	var err5 error
	go func() {
		defer wg.Done()
		endPlaceId, err5 = bsi.stationService.Get().QueryForStationId(ctx, info.EndPlace, token)
	}()

	wg.Wait()

	if err1 != nil {
		return util.TravelResult{}, err1
	}
	if err2 != nil {
		return util.TravelResult{}, err2
	}
	if err3 != nil {
		return util.TravelResult{}, err3
	}
	if err4 != nil {
		return util.TravelResult{}, err4
	}
	if err5 != nil {
		return util.TravelResult{}, err5
	}

	var indexStart int
	var indexEnd int

	for idx, val := range route.Stations {
		if val == startingPlaceId {
			indexStart = idx
			if indexEnd != 0 {
				break
			}
		}
		if val == endPlaceId {
			indexEnd = idx
			if indexStart != 0 {
				break
			}
		}
	}

	distance := float32(route.Distances[indexEnd] - route.Distances[indexStart])
	priceForEconomyClass := distance * priceConfig.BasicPriceRate
	priceForComfortClass := distance * priceConfig.FirstClassPriceRate

	return util.TravelResult{
		TrainType: train,
		Prices: map[string]float32{
			"EconomyClass": priceForEconomyClass,
			"ComfortClass": priceForComfortClass,
		},
		Percent: 1.0,
	}, nil

}

func (bsi *basicService) QueryForStationId(ctx context.Context, name, token string) (string, error) {
	err := util.Authenticate(token)
	if err != nil {
		return "", err
	}

	stationId, err := bsi.stationService.Get().QueryForStationId(ctx, name, token)
	if err != nil {
		return "", err
	}

	return stationId, nil
}
