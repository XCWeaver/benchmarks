package services

import (
	"context"
	"fmt"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type FoodMapService interface {
	GetAllFoodStores(ctx context.Context, token string) ([]util.Store, error)
	GetFoodStoresOfStation(ctx context.Context, stationId, token string) ([]util.Store, error)
	GetFoodStoresByStationIds(ctx context.Context, stationIds []string, token string) ([]util.Store, error)
	GetAllTrainFood(ctx context.Context, token string) ([]util.Food, error)
	GetTrainFoodOfTrip(ctx context.Context, tripId, token string) ([]util.Food, error)
}

type foodMapService struct {
	weaver.Implements[FoodMapService]
	//Mongo
	db components.NoSQLDatabase
}

func (fsi *foodMapService) GetAllFoodStores(ctx context.Context, token string) ([]util.Store, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("stores")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var stores []util.Store
	err = result.All(&stores)
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (fsi *foodMapService) GetFoodStoresOfStation(ctx context.Context, stationId, token string) ([]util.Store, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("stores")

	query := fmt.Sprintf(`{"StationId": %s }`, stationId)

	result, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var stores []util.Store
	err = result.All(&stores)
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (fsi *foodMapService) GetFoodStoresByStationIds(ctx context.Context, stationIds []string, token string) ([]util.Store, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("stores")

	query := fmt.Sprintf(`{"StationId": {"$in": %v} }`, stationIds)

	result, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var stores []util.Store
	err = result.All(&stores)
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (fsi *foodMapService) GetAllTrainFood(ctx context.Context, token string) ([]util.Food, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("trainfoods")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var foods []util.Food
	err = result.All(&foods)
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (fsi *foodMapService) GetTrainFoodOfTrip(ctx context.Context, tripId, token string) ([]util.Food, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("trainfoods")

	query := fmt.Sprintf(`{"TripId": %s }`, tripId)

	result, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var foods []util.Food
	err = result.All(&foods)
	if err != nil {
		return nil, err
	}

	return foods, nil
}
