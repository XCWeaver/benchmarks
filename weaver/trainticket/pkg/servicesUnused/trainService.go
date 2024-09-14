package services

import (
	"context"
	"errors"
	"log"

	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrainService interface {
	Create(ctx context.Context, train model.Train, token string) (model.Train, error)
	Update(ctx context.Context, train model.Train, token string) (model.Train, error)
	Delete(ctx context.Context, Id, token string) (string, error)
	Query(ctx context.Context, token string) ([]model.Train, error)
	Retrieve(ctx context.Context, Id, token string) (model.Train, error)
}

/*type trainServiceOptions struct {
	MongoAddr string `toml:"mongodb_address"`
	MongoPort string `toml:"mongodb_port"`
}*/

type trainService struct {
	weaver.Implements[TrainService]
	//weaver.WithConfig[trainServiceOptions]
	clientOptions *options.ClientOptions
	roles         []string
}

func (tsi *trainService) Init(ctx context.Context) error {
	//logger := tsi.Logger(ctx)

	//tsi.clientOptions = options.Client().ApplyURI("mongodb://" + tsi.Config().MongoAddr + ":" + tsi.Config().MongoPort + "/?directConnection=true")

	tsi.roles = append(tsi.roles, "role1")
	tsi.roles = append(tsi.roles, "role2")
	tsi.roles = append(tsi.roles, "role3")

	//logger.Info("train service running!", "mongodb_addr", tsi.Config().MongoAddr, "mongodb_port", tsi.Config().MongoPort)
	return nil
}

func (tsi *trainService) Create(ctx context.Context, train model.Train, token string) (model.Train, error) {
	logger := tsi.Logger(ctx)
	logger.Info("entering Create", "trainId", train.Id)

	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return model.Train{}, err
	}

	if train.Name == "" {
		return model.Train{}, errors.New("Train name not specified!")
	}

	client, err := mongo.Connect(ctx, tsi.clientOptions)
	if err != nil {
		return model.Train{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("trainType")
	filter := bson.D{{"name", train.Name}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == nil {
		return model.Train{}, errors.New("Train type already exists!")
	} else if result.Err() != mongo.ErrNoDocuments && result.Err() != nil {
		return model.Train{}, result.Err()
	}

	train.Id = uuid.New().String()
	res, err := collection.InsertOne(ctx, train)
	if err != nil {
		return model.Train{}, err
	}
	logger.Debug("train seccessfully created!", "objectid", res.InsertedID, "trainId", train.Id, "economyClass", train.EconomyClass,
		"confortClass", train.ComfortClass, "avgSpeed", train.AvgSpeed)

	return train, nil
}

func (tsi *trainService) Update(ctx context.Context, train model.Train, token string) (model.Train, error) {
	logger := tsi.Logger(ctx)
	logger.Info("entering Update", "trainId", train.Id)

	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return model.Train{}, err
	}

	client, err := mongo.Connect(ctx, tsi.clientOptions)
	if err != nil {
		return model.Train{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("trainType")
	filter := bson.D{{"id", train.Id}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return model.Train{}, errors.New("Train not found!")
	} else if result.Err() != nil {
		return model.Train{}, result.Err()
	}

	update := bson.D{{"$set", bson.D{{"name", train.Name}, {"economyClass", train.EconomyClass}, {"comfortClass", train.ComfortClass}, {"avgSpeed", train.AvgSpeed}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.Train{}, err
	}

	logger.Debug("train seccessfully updated!", "trainId", train.Id, "economyClass", train.EconomyClass,
		"confortClass", train.ComfortClass, "avgSpeed", train.AvgSpeed)

	return train, nil
}

func (tsi *trainService) Delete(ctx context.Context, Id, token string) (string, error) {
	logger := tsi.Logger(ctx)
	logger.Info("entering Delete", "trainId", Id)

	err := util.Authenticate(token, tsi.roles[0])
	if err != nil {
		return "", err
	}

	client, err := mongo.Connect(ctx, tsi.clientOptions)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("trainType")
	filter := bson.D{{"id", Id}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return "", errors.New("Train not found!")
	} else if result.Err() != nil {
		return "", result.Err()
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}

	return "Train type removed successfully", nil
}

func (tsi *trainService) Query(ctx context.Context, token string) ([]model.Train, error) {
	logger := tsi.Logger(ctx)
	logger.Info("entering Query")

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, tsi.clientOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("trainType")
	filter := bson.D{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var trains []model.Train
	if err = cursor.All(ctx, &trains); err != nil {
		return nil, err
	}
	logger.Debug("Get all users executed successfully!", "trains", trains)
	return trains, nil
}

func (tsi *trainService) Retrieve(ctx context.Context, Id, token string) (model.Train, error) {

	logger := tsi.Logger(ctx)
	logger.Info("entering Retrieve", "trainId", Id)

	err := util.Authenticate(token)

	if err != nil {
		return model.Train{}, err
	}

	client, err := mongo.Connect(ctx, tsi.clientOptions)
	if err != nil {
		return model.Train{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("trainType")
	filter := bson.D{{"id", Id}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return model.Train{}, errors.New("Train not found!")
	} else if result.Err() != nil {
		return model.Train{}, result.Err()
	}

	var existingTrain model.Train

	result.Decode(&existingTrain)

	logger.Debug("Train found!", "trainId", existingTrain.Id, "economyClass", existingTrain.EconomyClass, "comfortClass", existingTrain.ComfortClass, "avgSpeed", existingTrain.AvgSpeed)

	return existingTrain, nil
}
