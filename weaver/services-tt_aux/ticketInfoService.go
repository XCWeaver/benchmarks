package services

import (
	"context"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type TicketInfoService interface {
	QueryForTravel(ctx context.Context, info util.Travel, token string) (util.TravelResult, error)
	QueryForStationId(ctx context.Context, name, token string) (string, error)
}

type ticketInfoService struct {
	weaver.Implements[TicketInfoService]
	basicService weaver.Ref[BasicService]
}

func (tisi *ticketInfoService) QueryForTravel(ctx context.Context, info util.Travel, token string) (util.TravelResult, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.TravelResult{}, err
	}

	tr, err := tisi.basicService.Get().QueryForTravel(ctx, info, token)
	if err != nil {
		return util.TravelResult{}, err
	}

	return tr, nil
}

func (tisi *ticketInfoService) QueryForStationId(ctx context.Context, name, token string) (string, error) {
	err := util.Authenticate(token)
	if err != nil {
		return "", err
	}

	stationId, err := tisi.basicService.Get().QueryForStationId(ctx, name, token)
	if err != nil {
		return "", err
	}

	return stationId, nil
}
