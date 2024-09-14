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

type RouteService interface {
	CreateAndModifyRoute(ctx context.Context, route util.Route, token string) (util.Route, error)
	DeleteRoute(ctx context.Context, routeId, token string) (string, error)
	QueryById(ctx context.Context, routeId, token string) (util.Route, error)
	QueryAll(ctx context.Context, token string) ([]util.Route, error)
	QueryByStartAndTerminal(ctx context.Context, startId, terminalId, token string) ([]util.Route, error)
}

type routeService struct {
	weaver.Implements[RouteService]
	db    components.NoSQLDatabase
	roles []string
}

func (rsi *routeService) CreateAndModifyRoute(ctx context.Context, providedRoute util.RouteRequest, token string) (util.Route, error) {
	err := util.Authenticate(token, rsi.roles[0])

	if err != nil {
		return util.Route{}, err
	}

	stations := providedRoute.Stations
	distances := providedRoute.Distances

	if len(stations) != len(distances) {
		return util.Route{}, errors.New("Number of distances doesn't match that of stations.")
	}

	collection := rsi.db.GetDatabase("ts").GetCollection("routes")

	//* Insert
	if providedRoute.Id == "" || len(providedRoute.Id) < 10 {

		route := util.Route{
			Id:                uuid.New().String(),
			StartStationId:    providedRoute.StartStation,
			TerminalStationId: providedRoute.EndStation,
			Stations:          stations,
			Distances:         distances,
		}

		err = collection.InsertOne(route)
		if err != nil {
			return util.Route{}, err
		}

		return route, nil
	}

	//* Update

	query := fmt.Sprintf(`{"Id": %s }`, providedRoute.Id)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.Route{}, err
	}

	var existingRoute util.Route
	res.Decode(&existingRoute)

	if existingRoute.Id == "" {
		return util.Route{}, errors.New("Could not find route!")
	}

	route := util.Route{
		Id:                providedRoute.Id,
		StartStationId:    providedRoute.StartStation,
		TerminalStationId: providedRoute.EndStation,
		Stations:          stations,
		Distances:         distances,
	}

	err = collection.ReplaceOne(query, route)
	if err != nil {
		return util.Route{}, err
	}

	return route, nil
}

func (rsi *routeService) DeleteRoute(ctx context.Context, routeId, token string) (string, error) {
	err := util.Authenticate(token)

	if err != nil {
		return "", err
	}

	collection := rsi.db.GetDatabase("ts").GetCollection("routes")
	query := fmt.Sprintf(`{"Id": %s }`, routeId)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Route removed successfully", nil
}

func (rsi *routeService) QueryById(ctx context.Context, routeId, token string) (util.Route, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Route{}, err
	}

	collection := rsi.db.GetDatabase("ts").GetCollection("routes")

	query := fmt.Sprintf(`{"Id": %s}`, routeId)

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Route{}, err
	}

	var route util.Route
	err = result.Decode(&route)
	if err != nil {
		return util.Route{}, err
	}

	return route, nil
}

func (rsi *routeService) QueryAll(ctx context.Context, token string) ([]util.Route, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := rsi.db.GetDatabase("ts").GetCollection("routes")

	result, err := collection.FindMany("") //TODO verify this query-string works!

	if err != nil {
		return nil, err
	}

	var routes []util.Route
	err = result.All(&routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}

func (rsi *routeService) QueryByStartAndTerminal(ctx context.Context, startId, terminalId, token string) ([]util.Route, error) {

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := rsi.db.GetDatabase("ts").GetCollection("routes")

	query := fmt.Sprintf(`{"StartStationId": %s, "TerminalStationId": %s}`, startId, terminalId)

	result, err := collection.FindMany(query) //TODO verify this query-string works!

	if err != nil {
		return nil, err
	}

	var routes []util.Route
	err = result.All(&routes)
	if err != nil {
		return nil, err
	}

	var filteredRoutes []util.Route

	for _, route := range routes {

		foundStart := false
		for _, station := range route.Stations {

			if startId == station {
				foundStart = true
				continue //to find the end
			}

			if terminalId == station {

				if foundStart == false {
					break // route discarded cause startIdx > endIdx or startIdx not in stations
				}

				//cause start should have already been found at this point
				filteredRoutes = append(filteredRoutes, route)
				break
			}
		}
	}

	return filteredRoutes, nil
}
