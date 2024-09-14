package services

import (
	"context"
	"errors"
	"fmt"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type StationService interface {
	Query(ctx context.Context, token string) ([]util.Station, error)
	Create(ctx context.Context, station util.Station, token string) (util.Station, error)
	Update(ctx context.Context, station util.Station, token string) (util.Station, error)
	Delete(ctx context.Context, station util.Station, token string) (string, error)
	QueryForStationId(ctx context.Context, stationName, token string) (string, error)
	QueryForIdBatch(ctx context.Context, stationNameList []string, token string) ([]string, error)
	QueryById(ctx context.Context, stationId, token string) (string, error)
	QueryForNameBatch(ctx context.Context, stationIdList []string, token string) ([]string, error)
}

type stationService struct {
	weaver.Implements[StationService]
	db    components.NoSQLDatabase
	roles []string
}

func (ssi *stationService) Query(ctx context.Context, token string) ([]util.Station, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("stations")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var stations []util.Station
	err = result.All(&stations)
	if err != nil {
		return nil, err
	}

	return stations, nil
}

func (ssi *stationService) Create(ctx context.Context, station util.Station, token string) (util.Station, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return util.Station{}, err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("station")
	query := fmt.Sprintf(`{"Id": %s }`, station.Id)

	res, err := collection.FindOne(query)
	if err == nil {
		var oldStation util.Station
		res.Decode(&oldStation)

		if oldStation.Id != "" {
			return util.Station{}, errors.New("util.Station already exists!")
		}
	}

	err = collection.InsertOne(station)
	if err != nil {
		return util.Station{}, err
	}

	return station, nil
}

func (ssi *stationService) Update(ctx context.Context, station util.Station, token string) (util.Station, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return util.Station{}, err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("station")
	query := fmt.Sprintf(`{"Id": %s }`, station.Id)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.Station{}, err
	}

	var oldStation util.Station
	res.Decode(&oldStation)

	if oldStation.Id == "" {
		return util.Station{}, errors.New("util.Station does not exist!")
	}

	station.Id = oldStation.Id

	err = collection.ReplaceOne(query, station)
	if err != nil {
		return util.Station{}, err
	}

	return station, nil
}

func (ssi *stationService) Delete(ctx context.Context, station util.Station, token string) (string, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return "", err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("station")
	query := fmt.Sprintf(`{"Id": %s }`, station.Id)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Station removed successfully", nil
}

func (ssi *stationService) QueryForStationId(ctx context.Context, stationName, token string) (string, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return "", err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("station")
	query := fmt.Sprintf(`{"Name": %s }`, stationName)

	res, err := collection.FindOne(query)
	if err != nil {
		return "", err
	}

	var station util.Station
	res.Decode(&station)

	return station.Id, nil
}

func (ssi *stationService) QueryForIdBatch(ctx context.Context, stationNameList []string, token string) ([]string, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return nil, err
	}

	var query string
	collection := ssi.db.GetDatabase("ts").GetCollection("station")

	_, err = collection.FindOne(query)
	if err != nil {
		return nil, err
	}

	var ids []string
	var station util.Station

	for _, name := range stationNameList {
		query = fmt.Sprintf(`{"Name": %s`, name)
		res, err := collection.FindOne(query)
		if err == nil {
			res.Decode(&station)
			ids = append(ids, station.Id)
		}
	}

	return ids, nil
}

func (ssi *stationService) QueryById(ctx context.Context, stationId, token string) (string, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return "", err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("station")
	query := fmt.Sprintf(`{"Id": %s }`, stationId)

	res, err := collection.FindOne(query)
	if err != nil {
		return "", err
	}

	var station util.Station
	res.Decode(&station)

	return station.Name, nil
}

func (ssi *stationService) QueryForNameBatch(ctx context.Context, stationIdList []string, token string) ([]string, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return nil, err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("station")

	var names []string
	var query string
	var station util.Station

	for _, sId := range stationIdList {
		query = fmt.Sprintf(`{"Id": %s`, sId)
		res, err := collection.FindOne(query)
		if err == nil {
			res.Decode(&station)
			names = append(names, station.Name)
		}
	}

	return names, nil
}
