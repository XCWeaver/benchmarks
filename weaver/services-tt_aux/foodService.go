package services

//! **********************************************************************************************************************************************
//! **                                                                                                                                          **
//! **          This service  only communicates with Travel-service; Travel-2-service is not used. This might need an update.		            **
//! **                                                                                                                                          **
//! **********************************************************************************************************************************************

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type FoodService interface {
	FindAllFoodOrder(ctx context.Context, token string) ([]util.FoodOrder, error)
	CreateFoodOrder(ctx context.Context, foodOrder util.FoodOrder, token string) (util.FoodOrder, error)
	UpdateFoodOrder(ctx context.Context, foodOrder util.FoodOrder, token string) (util.FoodOrder, error)
	DeleteFoodOrder(ctx context.Context, orderId, token string) (string, error)
	FindFoodOrderByOrderId(ctx context.Context, orderId, token string) ([]util.FoodOrder, error)
	GetAllFood(ctx context.Context, starStation, endStation, tripId, token string) ([]util.Food, map[string]util.Store, error)
}

type foodService struct {
	weaver.Implements[FoodService]
	deliveryQueue components.Queue
	//Mongo
	db             components.NoSQLDatabase
	foodMapService weaver.Ref[FoodMapService]
	travelService  weaver.Ref[TravelService]
	stationService weaver.Ref[StationService]
	roles          []string
}

func (fsi *foodService) FindAllFoodOrder(ctx context.Context, token string) ([]util.FoodOrder, error) {

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("foodOrder")
	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var fos []util.FoodOrder
	err = result.All(&fos)
	if err != nil {
		return nil, err
	}
	return fos, nil
}

func (fsi *foodService) CreateFoodOrder(ctx context.Context, foodOrder util.FoodOrder, token string) (util.FoodOrder, error) {
	err := util.Authenticate(token, fsi.roles...)
	if err != nil {
		return util.FoodOrder{}, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("foodOrder")

	foodOrder.Id = uuid.New().String()

	if foodOrder.FoodType != 2 {
		foodOrder.StoreName = ""
		foodOrder.StationName = ""
	}

	err = collection.InsertOne(foodOrder)
	if err != nil {
		return util.FoodOrder{}, err
	}

	delivery := util.Delivery{FoodName: foodOrder.FoodName, ID: foodOrder.Id, StationName: foodOrder.StationName, StoreName: foodOrder.StoreName}
	delivery_bytes, err := json.Marshal(delivery)
	if err != nil {
		return foodOrder, err
	}
	err = fsi.deliveryQueue.Send(ctx, delivery_bytes)
	if err != nil {
		return foodOrder, err
	}
	return foodOrder, nil
}

func (fsi *foodService) UpdateFoodOrder(ctx context.Context, foodOrder util.FoodOrder, token string) (util.FoodOrder, error) {
	err := util.Authenticate(token, fsi.roles...)
	if err != nil {
		return util.FoodOrder{}, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("foodOrder")

	query := fmt.Sprintf(`{"Id": %s }`, foodOrder.Id)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.FoodOrder{}, err
	}

	var existingFo util.FoodOrder
	res.Decode(&existingFo)

	if existingFo.Id == "" {
		return util.FoodOrder{}, errors.New("Could not update food order!")
	}

	existingFo.FoodType = foodOrder.FoodType
	if foodOrder.FoodType == 1 {
		existingFo.StationName = foodOrder.StationName
		existingFo.StoreName = foodOrder.StoreName
	}

	existingFo.FoodName = foodOrder.FoodName
	existingFo.Price = foodOrder.Price

	query = fmt.Sprintf(`{"Id": %s}`, existingFo.Id)

	err = collection.ReplaceOne(query, existingFo)
	if err != nil {
		return util.FoodOrder{}, err
	}

	return existingFo, nil
}

func (fsi *foodService) DeleteFoodOrder(ctx context.Context, orderId, token string) (string, error) {
	err := util.Authenticate(token, fsi.roles...)
	if err != nil {
		return "", err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("trainType")

	query := fmt.Sprintf(`{"OrderId": %s }`, orderId)

	err = collection.DeleteMany(query)
	if err != nil {
		return "", err
	}

	return "util.Food orders removed successfully", nil
}

func (fsi *foodService) FindFoodOrderByOrderId(ctx context.Context, orderId, token string) ([]util.FoodOrder, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := fsi.db.GetDatabase("ts").GetCollection("foodOrder")
	query := fmt.Sprintf(`{"OrderId": %s }`, orderId)

	res, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var fos []util.FoodOrder

	res.All(&fos)

	return fos, nil
}

func (fsi *foodService) GetAllFood(ctx context.Context, startStation, endStation, tripId, token string) ([]util.Food, map[string]util.Store, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, nil, err
	}

	if len(tripId) < 3 {
		return nil, nil, errors.New("Invalid TripID")
	}

	var err1, err2 error
	var trainFoodList []util.Food
	var route util.Route

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		trainFoodList, err1 = fsi.foodMapService.Get().GetTrainFoodOfTrip(ctx, tripId, token)
	}()

	go func() {
		defer wg.Done()
		route, err2 = fsi.travelService.Get().GetRouteByTripId(ctx, tripId, token)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, nil, err1
	}
	if err2 != nil {
		return nil, nil, err2
	}

	stations := route.Stations
	aggStations := stations
	if startStation != "" {

		startStationId, err := fsi.stationService.Get().QueryForStationId(ctx, startStation, token)
		if err != nil {
			return nil, nil, err
		}

		for idx, s := range stations {
			if s == startStationId {
				break
			}

			aggStations[idx] = aggStations[len(aggStations)-1]
			aggStations = aggStations[:len(aggStations)-1]
		}
	}

	if endStation != "" {
		endStationId, err := fsi.stationService.Get().QueryForStationId(ctx, endStation, token)
		if err != nil {
			return nil, nil, err
		}

		for idx := len(stations) - 1; idx > 0; idx-- {
			if stations[idx] == endStationId {
				break
			}

			aggStations[idx] = aggStations[len(aggStations)-1]
			aggStations = aggStations[:len(aggStations)-1]
		}
	}

	stores, err := fsi.foodMapService.Get().GetFoodStoresByStationIds(ctx, aggStations, token)
	if err != nil {
		return nil, nil, err
	}

	foodStoreMap := make(map[string]util.Store)

	for _, s := range aggStations {
		for _, fs := range stores {
			foodStoreMap[s] = fs
		}
	}

	return trainFoodList, foodStoreMap, nil
}
