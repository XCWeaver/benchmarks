package services

import (
	"context"
	"fmt"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type ConsignPriceService interface {
	GetPriceByWeightAndRegion(ctx context.Context, weight float32, isWithinRegion bool, token string) (float32, error)
	GetPriceInfo(ctx context.Context, token string) (string, error)
	GetPriceConfig(ctx context.Context, token string) (util.ConsignPriceConfig, error)
	ModifyPriceConfig(ctx context.Context, priceConfig util.ConsignPriceConfig, token string) (util.ConsignPriceConfig, error)
}

type consignPriceService struct {
	weaver.Implements[ConsignPriceService]
	//Mongo
	db    components.NoSQLDatabase
	roles []string
}

func (cpsi *consignPriceService) GetPriceByWeightAndRegion(ctx context.Context, weight float32, isWithinRegion bool, token string) (float32, error) {
	err := util.Authenticate(token, cpsi.roles...)
	if err != nil {
		return 0.0, err
	}

	collection := cpsi.db.GetDatabase("ts").GetCollection("consign_price")
	query := fmt.Sprintf(`{"Index": %d }`, 0)

	res, err := collection.FindOne(query)
	if err != nil {
		return 0.0, err
	}

	var priceConfig util.ConsignPriceConfig
	res.Decode(&priceConfig)

	var price float32

	if weight < priceConfig.InitialWeight {
		price = priceConfig.InitialPrice
	} else {
		extraWeight := weight - priceConfig.InitialWeight

		if isWithinRegion {
			price = priceConfig.InitialWeight + extraWeight*priceConfig.WithinPrice
		} else {
			price = priceConfig.InitialWeight + extraWeight*priceConfig.BeyondPrice
		}
	}

	return price, nil
}

func (cpsi *consignPriceService) GetPriceInfo(ctx context.Context, token string) (string, error) {
	err := util.Authenticate(token, cpsi.roles...)
	if err != nil {
		return "", err
	}

	collection := cpsi.db.GetDatabase("ts").GetCollection("consign_price")
	query := fmt.Sprintf(`{"Index": %d }`, 0)

	res, err := collection.FindOne(query)
	if err != nil {
		return "", err
	}

	var priceConfig util.ConsignPriceConfig
	res.Decode(&priceConfig)

	info := fmt.Sprintf("The price of weight within %.2f is %.2f", priceConfig.InitialWeight, priceConfig.InitialPrice)

	return info, nil

}

func (cpsi *consignPriceService) GetPriceConfig(ctx context.Context, token string) (util.ConsignPriceConfig, error) {
	err := util.Authenticate(token, cpsi.roles...)
	if err != nil {
		return util.ConsignPriceConfig{}, err
	}

	collection := cpsi.db.GetDatabase("ts").GetCollection("consign_price")
	query := fmt.Sprintf(`{"Index": %d }`, 0)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.ConsignPriceConfig{}, err
	}

	var priceConfig util.ConsignPriceConfig
	res.Decode(&priceConfig)

	return priceConfig, nil
}

func (cpsi *consignPriceService) ModifyPriceConfig(ctx context.Context, priceConfig util.ConsignPriceConfig, token string) (util.ConsignPriceConfig, error) {
	err := util.Authenticate(token, cpsi.roles...)
	if err != nil {
		return util.ConsignPriceConfig{}, err
	}

	collection := cpsi.db.GetDatabase("ts").GetCollection("consign_price")
	query := fmt.Sprintf(`{"Index": %d }`, 0)

	_, err = collection.FindOne(query)

	if err != nil {
		err = collection.InsertOne(priceConfig)
		if err != nil {
			return util.ConsignPriceConfig{}, err
		}

	} else {
		priceConfig.Index = 0
		err = collection.ReplaceOne(query, priceConfig)
		if err != nil {
			return util.ConsignPriceConfig{}, err
		}
	}

	return priceConfig, nil
}
