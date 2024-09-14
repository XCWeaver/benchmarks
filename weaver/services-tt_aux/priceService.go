package services

import (
	"context"
	"errors"
	"fmt"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type PriceService interface {
	Query(ctx context.Context, routeId, trainType, token string) (util.PriceConfig, error)
	QueryAll(ctx context.Context, token string) ([]util.PriceConfig, error)
	Create(ctx context.Context, pc util.PriceConfig, token string) (util.PriceConfig, error)
	Delete(ctx context.Context, pc util.PriceConfig, token string) (string, error)
	Update(ctx context.Context, pc util.PriceConfig, token string) (util.PriceConfig, error)
}

type priceService struct {
	weaver.Implements[PriceService]
	db    components.NoSQLDatabase
	roles []string
}

func (psi *priceService) Query(ctx context.Context, routeId, trainType, token string) (util.PriceConfig, error) {
	err := util.Authenticate(token)

	if err != nil {
		return util.PriceConfig{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")

	query := fmt.Sprintf(`{"RouteId": %s, "TrainType": %s}`, routeId, trainType)

	result, err := collection.FindOne(query)

	if err != nil {
		return util.PriceConfig{}, err
	}

	var pc util.PriceConfig
	err = result.Decode(&pc)
	if err != nil {
		return util.PriceConfig{}, err
	}

	return pc, nil
}

func (psi *priceService) QueryAll(ctx context.Context, token string) ([]util.PriceConfig, error) {
	err := util.Authenticate(token)

	if err != nil {
		return nil, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")
	result, err := collection.FindMany("") //TODO verify this query-string works!

	if err != nil {
		return nil, err
	}

	var pcs []util.PriceConfig
	err = result.All(&pcs)
	if err != nil {
		return nil, err
	}

	return pcs, nil
}

func (psi *priceService) Create(ctx context.Context, pc util.PriceConfig, token string) (util.PriceConfig, error) {
	err := util.Authenticate(token, psi.roles...)

	if err != nil {
		return util.PriceConfig{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")

	query := fmt.Sprintf(`{"Id": %s }`, pc.Id)

	res, err := collection.FindOne(query)

	if err == nil {
		var oldPc util.PriceConfig

		res.Decode(&oldPc)

		if oldPc.Id != "" {
			return util.PriceConfig{}, errors.New("Price Config already exists!")
		}
	}

	pc.Id = uuid.New().String()
	err = collection.InsertOne(pc)
	if err != nil {
		return util.PriceConfig{}, err
	}

	return pc, nil
}

func (psi *priceService) Delete(ctx context.Context, pc util.PriceConfig, token string) (string, error) {
	err := util.Authenticate(token, psi.roles...)

	if err != nil {
		return "", err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")
	query := fmt.Sprintf(`{"Id": %s }`, pc.Id)

	res, err := collection.FindOne(query)

	if err != nil {
		return "", err
	}

	var existingPc util.PriceConfig
	res.Decode(&existingPc)
	if existingPc.Id == "" {
		return "", errors.New("Could not find price config!")
	}

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "Price Config removed successfully", nil
}

func (psi *priceService) Update(ctx context.Context, pc util.PriceConfig, token string) (util.PriceConfig, error) {
	err := util.Authenticate(token, psi.roles...)

	if err != nil {
		return util.PriceConfig{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")

	query := fmt.Sprintf(`{"Id": %s }`, pc.Id)

	res, err := collection.FindOne(query)

	if err != nil {
		return util.PriceConfig{}, err
	}

	var existingPc util.PriceConfig
	res.Decode(&existingPc)

	if existingPc.Id == "" {
		return util.PriceConfig{}, errors.New("Could not find price config!")
	}

	query = fmt.Sprintf(`{"Id": %s}`, pc.Id)
	update := fmt.Sprintf(`{"$set": {"BasicPriceRate": %d, "FirstClassPriceRate": %d, "RouteId": %s, "TrainType": %s}`, pc.BasicPriceRate, pc.FirstClassPriceRate, pc.RouteId, pc.TrainType)

	err = collection.UpdateOne(query, update)
	if err != nil {
		return util.PriceConfig{}, err
	}

	return pc, nil
}
