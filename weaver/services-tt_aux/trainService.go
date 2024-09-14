package services

import (
	"context"
	"errors"
	"fmt"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type TrainService interface {
	Create(ctx context.Context, train util.Train, token string) (util.Train, error)
	Update(ctx context.Context, train util.Train, token string) (util.Train, error)
	Delete(ctx context.Context, Id, token string) (string, error)
	Query(ctx context.Context, token string) ([]util.Train, error)
	Retrieve(ctx context.Context, Id, token string) (util.Train, error)
}

type trainService struct {
	weaver.Implements[TrainService]
	db    components.NoSQLDatabase
	roles []string
}

// * This endpoint expects the `train` object to have an ID
func (tsi *trainService) Create(ctx context.Context, train util.Train, token string) (util.Train, error) {

	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return util.Train{}, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trainType")

	query := fmt.Sprintf(`{"Id": %s }`, train.Id)

	res, err := collection.FindOne(query)

	if err == nil {
		var oldTrain util.Train

		res.Decode(&oldTrain)

		if oldTrain.Id != "" {
			return util.Train{}, errors.New("util.Train type already exists!")
		}
	}

	err = collection.InsertOne(train)
	if err != nil {
		return util.Train{}, err
	}

	return train, nil
}

func (tsi *trainService) Update(ctx context.Context, train util.Train, token string) (util.Train, error) {
	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return util.Train{}, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trainType")

	query := fmt.Sprintf(`{"Id": %s }`, train.Id)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.Train{}, err
	}

	var existingTrain util.Train
	res.Decode(&existingTrain)

	if existingTrain.Id == "" {
		return util.Train{}, errors.New("Could not update train type!")
	}

	query = fmt.Sprintf(`{"Id": %s}`, existingTrain.Id)
	update := fmt.Sprintf(`{"$set": {"EconomyClass": %d, "ComfortClass": %d, "AvgSpeed": %d}`, train.EconomyClass, train.ComfortClass, train.AvgSpeed)

	err = collection.UpdateOne(query, update)
	if err != nil {
		return util.Train{}, err
	}

	return train, nil
}

func (tsi *trainService) Delete(ctx context.Context, Id, token string) (string, error) {
	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return "", err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trainType")

	query := fmt.Sprintf(`{"Id": %s }`, Id)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Train type removed successfully", nil
}

func (tsi *trainService) Query(ctx context.Context, token string) ([]util.Train, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := tsi.db.GetDatabase("ts").GetCollection("trainType")
	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var trains []util.Train
	err = result.All(&trains)
	if err != nil {
		return nil, err
	}
	return trains, nil
}

func (tsi *trainService) Retrieve(ctx context.Context, Id, token string) (util.Train, error) {

	err := util.Authenticate(token)

	if err != nil {
		return util.Train{}, err
	}

	query := fmt.Sprintf(`{"Id": %s }`, Id)
	collection := tsi.db.GetDatabase("ts").GetCollection("trainType")

	res, err := collection.FindOne(query)
	if err != nil {
		return util.Train{}, err
	}

	var existingTrain util.Train

	res.Decode(&existingTrain)

	return existingTrain, nil
}
