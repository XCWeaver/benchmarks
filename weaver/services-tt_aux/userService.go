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

type UserService interface {
	GetAllUsers(ctx context.Context, token string) ([]util.User, error)
	GetUserById(ctx context.Context, userId, token string) (util.User, error)
	GetUserByUsername(ctx context.Context, username, token string) (util.User, error)
	DeleteUserById(ctx context.Context, userId, token string) (string, error)
	UpdateUser(ctx context.Context, userData util.User, token string) (util.User, error)
	RegisterUser(ctx context.Context, userData util.User, token string) (util.User, error)
}

type userService struct {
	weaver.Implements[UserService]
	db          components.NoSQLDatabase
	authService weaver.Ref[AuthService]
	roles       []string
}

func (usi *userService) GetAllUsers(ctx context.Context, token string) ([]util.User, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := usi.db.GetDatabase("ts-user-mongo").GetCollection("user")
	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var users []util.User
	err = result.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (usi *userService) GetUserById(ctx context.Context, userId, token string) (util.User, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.User{}, err
	}

	query := fmt.Sprintf(`{"UserId": %s }`, userId)
	collection := usi.db.GetDatabase("ts-user-mongo").GetCollection("user")

	res, err := collection.FindOne(query)
	if err != nil {
		return util.User{}, err
	}

	var existingUser util.User

	res.Decode(&existingUser)

	return existingUser, nil
}

func (usi *userService) GetUserByUsername(ctx context.Context, username, token string) (util.User, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.User{}, err
	}

	query := fmt.Sprintf(`{"Username": %s }`, username)
	collection := usi.db.GetDatabase("ts-user-mongo").GetCollection("user")

	res, err := collection.FindOne(query)
	if err != nil {
		return util.User{}, err
	}

	var existingUser util.User

	res.Decode(&existingUser)

	return existingUser, nil
}

func (usi *userService) DeleteUserById(ctx context.Context, userId, token string) (string, error) {
	err := util.Authenticate(token, usi.roles...)
	if err != nil {
		return "", err
	}

	collection := usi.db.GetDatabase("ts-user-mongo").GetCollection("user")

	query := fmt.Sprintf(`{"UserId": %s }`, userId)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.User removed successfully", nil
}

// * userData:util.User should not contain the userId
func (usi *userService) UpdateUser(ctx context.Context, userData util.User, token string) (util.User, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.User{}, err
	}

	collection := usi.db.GetDatabase("ts-user-mongo").GetCollection("user")
	query := fmt.Sprintf(`{"Username": %s }`, userData.Username)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.User{}, err
	}

	var existingUser util.User
	res.Decode(&existingUser)

	delQuery := fmt.Sprintf(`{"UserId": %s }`, existingUser.UserId)
	err = collection.DeleteOne(delQuery)
	if err != nil {
		return util.User{}, err
	}

	userData.UserId = existingUser.UserId
	err = collection.InsertOne(userData)
	if err != nil {
		return util.User{}, err
	}

	return userData, nil
}

// * userData:util.User should not contain the userId
func (usi *userService) RegisterUser(ctx context.Context, userData util.User, token string) (util.User, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.User{}, err
	}

	collection := usi.db.GetDatabase("ts-user-mongo").GetCollection("user")
	query := fmt.Sprintf(`{"Username": %s }`, userData.Username)

	res, err := collection.FindOne(query)
	if err == nil {
		var existingUser util.User
		res.Decode(&existingUser)
		if existingUser.UserId != "" {
			return util.User{}, errors.New("util.User with this username already registered")
		}
	}

	userData.UserId = uuid.New().String()
	err = collection.InsertOne(userData)
	if err != nil {
		return util.User{}, err
	}

	_, err = usi.authService.Get().CreateDefaultUser(ctx, userData, token)
	if err != nil {
		return util.User{}, err
	}

	return userData, nil
}
