package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type ConsignService interface {
	InsertConsign(ctx context.Context, consign util.Consign, token string) (util.Consign, error)
	UpdateConsign(ctx context.Context, consign util.Consign, token string) (util.Consign, error)
	FindByAccountId(ctx context.Context, accountId, token string) ([]util.Consign, error)
	FindByOrderId(ctx context.Context, orderId, token string) ([]util.Consign, error)
	FindByConsignee(ctx context.Context, consignee, token string) ([]util.Consign, error)
}

type consignService struct {
	weaver.Implements[ConsignService]
	consignPriceService weaver.Ref[ConsignPriceService]
	//Mongo
	db    components.NoSQLDatabase
	roles []string
}

func (csi *consignService) InsertConsign(ctx context.Context, consign util.Consign, token string) (util.Consign, error) {
	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return util.Consign{}, err
	}

	price, err := csi.consignPriceService.Get().GetPriceByWeightAndRegion(ctx, consign.Weight, consign.Within, token)
	if err != nil {
		return util.Consign{}, err
	}

	consign.Price = price
	consign.Id = uuid.New().String()

	collection := csi.db.GetDatabase("ts").GetCollection("consign_record")
	err = collection.InsertOne(consign)
	if err != nil {
		return util.Consign{}, err
	}

	return consign, nil
}

func (csi *consignService) UpdateConsign(ctx context.Context, consign util.Consign, token string) (util.Consign, error) {

	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return util.Consign{}, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("consign_record")

	query := fmt.Sprintf(`{"Id": %s }`, consign.Id)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.Consign{}, err
	}

	var originalConsign util.Consign
	res.Decode(originalConsign)

	if originalConsign.Weight != consign.Weight {
		retval, err := csi.consignPriceService.Get().GetPriceByWeightAndRegion(ctx, consign.Weight, consign.Within, token)
		if err != nil {
			return util.Consign{}, err
		}

		consign.Price = retval
	}

	err = collection.ReplaceOne(query, consign)
	if err != nil {
		return util.Consign{}, err
	}

	return consign, nil
}

func (csi *consignService) FindByAccountId(ctx context.Context, accountId, token string) ([]util.Consign, error) {
	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return nil, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("consign_record")
	query := fmt.Sprintf(`{"AccountId": %s}`, accountId)

	result, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var consigns []util.Consign
	err = result.All(&consigns)
	if err != nil {
		return nil, err
	}

	return consigns, nil
}

func (csi *consignService) FindByOrderId(ctx context.Context, orderId, token string) ([]util.Consign, error) {
	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return nil, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("consign_record")
	query := fmt.Sprintf(`{"OrderId": %s}`, orderId)

	result, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var consigns []util.Consign
	err = result.All(&consigns)
	if err != nil {
		return nil, err
	}

	return consigns, nil
}

func (csi *consignService) FindByConsignee(ctx context.Context, consignee, token string) ([]util.Consign, error) {
	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return nil, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("consign_record")
	query := fmt.Sprintf(`{"Consignee": %s}`, consignee)

	result, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var consigns []util.Consign
	err = result.All(&consigns)
	if err != nil {
		return nil, err
	}

	return consigns, nil
}
