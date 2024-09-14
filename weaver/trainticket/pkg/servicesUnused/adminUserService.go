package services

import (
	"context"
	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type AdminUserService interface {
	GetAllUsers(ctx context.Context, token string) ([]model.User, error)
	AddUser(ctx context.Context, user model.User, token string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User, token string) (model.User, error)
	DeleteUser(ctx context.Context, userId, token string) (string, error)
}

type adminUserService struct {
	weaver.Implements[AdminUserService]
	userService weaver.Ref[UserService]
	roles       []string
}

func (ausi *adminUserService) GetAllUsers(ctx context.Context, token string) ([]model.User, error) {

	err := util.Authenticate(token, ausi.roles...)
	if err != nil {
		return nil, err
	}

	return ausi.userService.Get().GetAllUsers(ctx, token)
}

func (ausi *adminUserService) AddUser(ctx context.Context, user model.User, token string) (model.User, error) {

	err := util.Authenticate(token, ausi.roles...)
	if err != nil {
		return model.User{}, err
	}

	return ausi.userService.Get().RegisterUser(ctx, user, token)
}

func (ausi *adminUserService) UpdateUser(ctx context.Context, user model.User, token string) (model.User, error) {

	err := util.Authenticate(token, ausi.roles...)
	if err != nil {
		return model.User{}, err
	}

	return ausi.userService.Get().UpdateUser(ctx, user, token)
}

func (ausi *adminUserService) DeleteUser(ctx context.Context, userId, token string) (string, error) {

	err := util.Authenticate(token, ausi.roles...)
	if err != nil {
		return "", err
	}

	return ausi.userService.Get().DeleteUserById(ctx, userId, token)
}
