package services

import (
	"context"
	"errors"
	"fmt"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type ConfigService interface {
	QueryAll(ctx context.Context) ([]util.Config, error)
	CreateConfig(ctx context.Context, info util.Config) (util.Config, error)
	UpdateConfig(ctx context.Context, info util.Config) (util.Config, error)
	DeleteConfig(ctx context.Context, configName string) (string, error)
	Retrieve(ctx context.Context, configName string) (util.Config, error)
}

type configService struct {
	weaver.Implements[ConfigService]
	//Mongo
	db components.NoSQLDatabase
}

func (csi *configService) QueryAll(ctx context.Context) ([]util.Config, error) {
	collection := csi.db.GetDatabase("ts").GetCollection("config")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var configs []util.Config
	err = result.All(&configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (csi *configService) CreateConfig(ctx context.Context, info util.Config) (util.Config, error) {
	collection := csi.db.GetDatabase("ts").GetCollection("config")

	query := fmt.Sprintf(`{"Name": %s }`, info.Name)

	result, err := collection.FindOne(query)
	if err == nil {

		var oldConfig util.Config
		err = result.Decode(&oldConfig)
		if err == nil {
			return util.Config{}, errors.New("A config with this name already exists!")
		}
	}

	err = collection.InsertOne(info)
	if err != nil {
		return util.Config{}, err
	}

	return info, nil
}

func (csi *configService) UpdateConfig(ctx context.Context, info util.Config) (util.Config, error) {
	collection := csi.db.GetDatabase("ts").GetCollection("config")

	query := fmt.Sprintf(`{"Name": %s }`, info.Name)

	_, err := collection.FindOne(query)
	if err != nil {
		return util.Config{}, err
	}

	err = collection.ReplaceOne(query, info)
	if err != nil {
		return util.Config{}, err
	}

	return info, nil
}

func (csi *configService) DeleteConfig(ctx context.Context, configName string) (string, error) {
	collection := csi.db.GetDatabase("ts").GetCollection("config")

	query := fmt.Sprintf(`{"Name": %s }`, configName)

	_, err := collection.FindOne(query)
	if err != nil {
		return "", err
	}

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Config deleted.", nil
}

func (csi *configService) Retrieve(ctx context.Context, configName string) (util.Config, error) {

	collection := csi.db.GetDatabase("ts").GetCollection("config")

	query := fmt.Sprintf(`{"Name": %s }`, configName)
	result, err := collection.FindOne(query)
	if err != nil {
		return util.Config{}, err
	}

	var config util.Config
	err = result.Decode(&config)
	if err != nil {
		return util.Config{}, err
	}

	return config, nil
}
