package services

import (
	"context"
	"log"

	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/TiagoMalhadas/xcweaver"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PriceService interface {
	//Query(ctx context.Context, routeId, trainType, token string) (util.PriceConfig, error)
	//QueryAll(ctx context.Context, token string) ([]util.PriceConfig, error)
	Create(ctx context.Context, pc model.PriceConfig, token string) (model.PriceConfig, error)
	//Delete(ctx context.Context, pc util.PriceConfig, token string) (string, error)
	//Update(ctx context.Context, pc util.PriceConfig, token string) (util.PriceConfig, error)
}

/*type priceServiceOptions struct {
	MongoAddr string `toml:"mongodb_address"`
	MongoPort string `toml:"mongodb_port"`
}*/

type priceService struct {
	xcweaver.Implements[PriceService]
	//xcweaver.WithConfig[priceServiceOptions]
	clientOptions *options.ClientOptions
	roles         []string
}

func (psi *priceService) Init(ctx context.Context) error {
	//logger := psi.Logger(ctx)

	//psi.clientOptions = options.Client().ApplyURI("mongodb://" + psi.Config().MongoAddr + ":" + psi.Config().MongoPort + "/?directConnection=true")

	psi.roles = append(psi.roles, "role1")
	psi.roles = append(psi.roles, "role2")
	psi.roles = append(psi.roles, "role3")

	//logger.Info("price service running!", "mongodb_addr", psi.Config().MongoAddr, "mongodb_port", psi.Config().MongoPort)
	return nil
}

/*func (psi *priceService) Query(ctx context.Context, routeId, trainType, token string) (model.PriceConfig, error) {
	err := util.Authenticate(token)

	if err != nil {
		return model.PriceConfig{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")

	query := fmt.Sprintf(`{"RouteId": %s, "TrainType": %s}`, routeId, trainType)

	result, err := collection.FindOne(query)

	if err != nil {
		return model.PriceConfig{}, err
	}

	var pc model.PriceConfig
	err = result.Decode(&pc)
	if err != nil {
		return model.PriceConfig{}, err
	}

	return pc, nil
}

func (psi *priceService) QueryAll(ctx context.Context, token string) ([]model.PriceConfig, error) {
	err := util.Authenticate(token)

	if err != nil {
		return nil, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")
	result, err := collection.FindMany("") //TODO verify this query-string works!

	if err != nil {
		return nil, err
	}

	var pcs []model.PriceConfig
	err = result.All(&pcs)
	if err != nil {
		return nil, err
	}

	return pcs, nil
}*/

func (psi *priceService) Create(ctx context.Context, pc model.PriceConfig, token string) (model.PriceConfig, error) {
	logger := psi.Logger(ctx)
	logger.Info("entering Create", "trainId", pc.Id)

	err := util.Authenticate(token, psi.roles...)

	if err != nil {
		return model.PriceConfig{}, err
	}

	client, err := mongo.Connect(ctx, psi.clientOptions)
	if err != nil {
		return model.PriceConfig{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("price_config")

	//create new
	if pc.Id == "" {
		pc.Id = uuid.New().String()
		result, err := collection.InsertOne(ctx, pc)
		if err != nil {
			return model.PriceConfig{}, err
		}
		logger.Debug("inserted new price config", "objectid", result.InsertedID, "priceConfigId", pc.Id, "trainType", pc.TrainType, "routeId", pc.RouteId,
			"basicPriceRate", pc.BasicPriceRate, "firstClasspriceRate", pc.FirstClassPriceRate)
	} else {
		//update
		filter := bson.D{{"id", pc.Id}}
		res := collection.FindOne(ctx, filter)
		if res.Err() == mongo.ErrNoDocuments {
			result, err := collection.InsertOne(ctx, pc)
			if err != nil {
				return model.PriceConfig{}, err
			}
			logger.Debug("inserted new price config", "objectid", result.InsertedID, "priceConfigId", pc.Id, "trainType", pc.TrainType, "routeId", pc.RouteId,
				"basicPriceRate", pc.BasicPriceRate, "firstClasspriceRate", pc.FirstClassPriceRate)
		} else if res.Err() != nil {
			return model.PriceConfig{}, err
		} else {
			update := bson.D{{"$set", bson.D{{"trainType", pc.TrainType}, {"routeId", pc.RouteId}, {"basicPriceRate", pc.BasicPriceRate}, {"firstClassPriceRate", pc.FirstClassPriceRate}}}}
			result, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				return model.PriceConfig{}, err
			}
			logger.Debug("price config updated", "objectid", result.UpsertedID, "priceConfigId", pc.Id, "trainType", pc.TrainType, "routeId", pc.RouteId,
				"basicPriceRate", pc.BasicPriceRate, "firstClasspriceRate", pc.FirstClassPriceRate)
		}
	}

	return pc, nil
}

/*func (psi *priceService) Delete(ctx context.Context, pc model.PriceConfig, token string) (string, error) {
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

	var existingPc model.PriceConfig
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

func (psi *priceService) Update(ctx context.Context, pc model.PriceConfig, token string) (model.PriceConfig, error) {
	err := util.Authenticate(token, psi.roles...)

	if err != nil {
		return model.PriceConfig{}, err
	}

	collection := psi.db.GetDatabase("ts").GetCollection("price_config")

	query := fmt.Sprintf(`{"Id": %s }`, pc.Id)

	res, err := collection.FindOne(query)

	if err != nil {
		return model.PriceConfig{}, err
	}

	var existingPc model.PriceConfig
	res.Decode(&existingPc)

	if existingPc.Id == "" {
		return model.PriceConfig{}, errors.New("Could not find price config!")
	}

	query = fmt.Sprintf(`{"Id": %s}`, pc.Id)
	update := fmt.Sprintf(`{"$set": {"BasicPriceRate": %d, "FirstClassPriceRate": %d, "RouteId": %s, "TrainType": %s}`, pc.BasicPriceRate, pc.FirstClassPriceRate, pc.RouteId, pc.TrainType)

	err = collection.UpdateOne(query, update)
	if err != nil {
		return model.PriceConfig{}, err
	}

	return pc, nil
}*/
