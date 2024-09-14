package services

import (
	"context"
	"errors"
	"log"

	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/TiagoMalhadas/xcweaver"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService interface {
	GetAllUsers(ctx context.Context, token string) ([]model.User, error)
	GetUserById(ctx context.Context, userId, token string) (model.User, error)
	GetUserByUsername(ctx context.Context, username, token string) (model.User, error)
	DeleteUserById(ctx context.Context, userId, token string) (string, error)
	UpdateUser(ctx context.Context, userData model.User, token string) (model.User, error)
	RegisterUser(ctx context.Context, userData model.User, token string) (model.User, error)
}

/*type userServiceOptions struct {
	MongoAddr string `toml:"mongodb_address"`
	MongoPort string `toml:"mongodb_port"`
}*/

type userService struct {
	xcweaver.Implements[UserService]
	//xcweaver.WithConfig[userServiceOptions]
	clientOptions *options.ClientOptions
	authService   xcweaver.Ref[AuthService]
	roles         []string
}

func (u *userService) Init(ctx context.Context) error {
	//logger := u.Logger(ctx)

	//u.clientOptions = options.Client().ApplyURI("mongodb://" + u.Config().MongoAddr + ":" + u.Config().MongoPort + "/?directConnection=true")

	//logger.Info("user service running!", "mongodb_addr", u.Config().MongoAddr, "mongodb_port", u.Config().MongoPort)
	return nil
}

func (usi *userService) GetAllUsers(ctx context.Context, token string) ([]model.User, error) {
	logger := usi.Logger(ctx)
	logger.Info("entering GetAllUsers")

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, usi.clientOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-user-mongo").Collection("user")
	filter := bson.D{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var users []model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	logger.Debug("Get all users executed successfully!", "users", users)

	return users, nil
}

func (usi *userService) GetUserById(ctx context.Context, userId, token string) (model.User, error) {
	logger := usi.Logger(ctx)
	logger.Info("entering GetUserById", "userId", userId)

	err := util.Authenticate(token)
	if err != nil {
		return model.User{}, err
	}

	client, err := mongo.Connect(ctx, usi.clientOptions)
	if err != nil {
		return model.User{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-user-mongo").Collection("user")
	filter := bson.D{{"user_id", userId}}

	var existingUser model.User
	err = collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, err
	}

	logger.Debug("user successfully found!", "username", existingUser.Username, "password", existingUser.Password, "role", existingUser.Role, "userId", existingUser.UserId,
		"email", existingUser.Email, "documentType", existingUser.DocumentType, "documentNumber", existingUser.DocumentNumber, "gender", existingUser.Gender)

	return existingUser, nil
}

func (usi *userService) GetUserByUsername(ctx context.Context, username, token string) (model.User, error) {
	logger := usi.Logger(ctx)
	logger.Info("entering GetUserByUsername", "username", username)

	err := util.Authenticate(token)
	if err != nil {
		return model.User{}, err
	}

	client, err := mongo.Connect(ctx, usi.clientOptions)
	if err != nil {
		return model.User{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-user-mongo").Collection("user")
	filter := bson.D{{"username", username}}

	var existingUser model.User
	err = collection.FindOne(ctx, filter).Decode(&existingUser)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, err
	}

	logger.Debug("user successfully found!", "username", existingUser.Username, "password", existingUser.Password, "role", existingUser.Role, "userId", existingUser.UserId,
		"email", existingUser.Email, "documentType", existingUser.DocumentType, "documentNumber", existingUser.DocumentNumber, "gender", existingUser.Gender)

	return existingUser, nil
}

func (usi *userService) DeleteUserById(ctx context.Context, userId, token string) (string, error) {
	logger := usi.Logger(ctx)
	logger.Info("entering DeleteUserById", "userId", userId)

	err := util.Authenticate(token, usi.roles...)
	if err != nil {
		return "", err
	}

	client, err := mongo.Connect(ctx, usi.clientOptions)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-user-mongo").Collection("user")
	filter := bson.D{{"user_id", userId}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return "", errors.New("User not found!")
	} else if result.Err() != nil {
		return "", result.Err()
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}

	return "User removed successfully", nil
}

// * userData:model.User should not contain the userId
func (usi *userService) UpdateUser(ctx context.Context, userData model.User, token string) (model.User, error) {
	logger := usi.Logger(ctx)
	logger.Info("entering UpdateUser", "username", userData.Username)

	err := util.Authenticate(token)
	if err != nil {
		return model.User{}, err
	}

	client, err := mongo.Connect(ctx, usi.clientOptions)
	if err != nil {
		return model.User{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-user-mongo").Collection("user")
	filter := bson.D{{"username", userData.Username}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return model.User{}, errors.New("User not found!")
	} else if result.Err() != nil {
		return model.User{}, result.Err()
	}

	update := bson.D{{"$set", bson.D{{"password", userData.Password}, {"role", userData.Role}, {"email", userData.Email},
		{"documentType", userData.DocumentType}, {"documentNumber", userData.DocumentNumber}, {"gender", userData.Gender}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.User{}, err
	}

	logger.Debug("user successfully updated!", "username", userData.Username, "password", userData.Password, "role", userData.Role, "userId", userData.UserId,
		"email", userData.Email, "documentType", userData.DocumentType, "documentNumber", userData.DocumentNumber, "gender", userData.Gender)

	return userData, nil
}

// * userData:model.User should not contain the userId
func (u *userService) RegisterUser(ctx context.Context, userData model.User, token string) (model.User, error) {
	logger := u.Logger(ctx)
	logger.Info("entering RegisterUser", "username", userData.Username)

	/*err := util.Authenticate(token)
	if err != nil {
		return model.User{}, err
	}*/

	client, err := mongo.Connect(ctx, u.clientOptions)
	if err != nil {
		return model.User{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-user-mongo").Collection("user")
	filter := bson.D{{"username", userData.Username}}

	res := collection.FindOne(ctx, filter)
	if res.Err() == nil {
		return model.User{}, errors.New("User with this username already registered!")
	} else if res.Err() != mongo.ErrNoDocuments && res.Err() != nil {
		return model.User{}, res.Err()
	}

	userData.UserId = uuid.New().String()
	result, err := collection.InsertOne(ctx, userData)
	if err != nil {
		return model.User{}, err
	}
	logger.Debug("inserted user", "objectid", result.InsertedID)

	_, err = u.authService.Get().CreateDefaultUser(ctx, userData, token)
	if err != nil {
		return model.User{}, err
	}

	logger.Debug("Default user created successfully", "username", userData.Username, "password", userData.Password, "role", userData.Role, "userId", userData.UserId, "email",
		userData.Email, "documentType", userData.DocumentType, "documentNumber", userData.DocumentNumber, "gender", userData.Gender)

	return userData, nil
}
