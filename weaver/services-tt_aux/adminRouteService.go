package services

import (
	"context"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type AdminRouteService interface {
	GetAllRoutes(ctx context.Context, token string) ([]util.Route, error)
	AddRoute(ctx context.Context, route util.RouteRequest, token string) (util.Route, error)
	DeleteRoute(ctx context.Context, routeId, token string) (string, error)
}

type adminRouteService struct {
	weaver.Implements[AdminRouteService]
	routeService weaver.Ref[RouteService]
	roles        []string
}

func (arsi *adminRouteService) GetAllRoutes(ctx context.Context, token string) ([]util.Route, error) {

	err := util.Authenticate(token, arsi.roles...)
	if err != nil {
		return nil, err
	}

	return arsi.routeService.Get().QueryAll(ctx, token)
}

func (arsi *adminRouteService) AddRoute(ctx context.Context, routereq util.RouteRequest, token string) (util.Route, error) {

	err := util.Authenticate(token, arsi.roles...)
	if err != nil {
		return util.Route{}, err
	}
	route := util.Route{Id: routereq.Id, StartStationId: routereq.StartStation, TerminalStationId: routereq.EndStation, Stations: routereq.Stations, Distances: routereq.Distances}
	return arsi.routeService.Get().CreateAndModifyRoute(ctx, route, token)
}

func (arsi *adminRouteService) DeleteRoute(ctx context.Context, routeId, token string) (string, error) {
	err := util.Authenticate(token, arsi.roles...)
	if err != nil {
		return "", err
	}

	return arsi.routeService.Get().DeleteRoute(ctx, routeId, token)
}
