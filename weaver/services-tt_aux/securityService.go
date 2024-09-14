package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type SecurityService interface {
	FindAllSecurityConfigs(ctx context.Context, token string) ([]util.SecurityConfig, error)
	Create(ctx context.Context, name, value, description, token string) (util.SecurityConfig, error)
	Update(ctx context.Context, id, name, value, description, token string) (util.SecurityConfig, error)
	Delete(ctx context.Context, id, token string) (string, error)
	Check(ctx context.Context, accountId, token string) (string, error)
}

type securityService struct {
	weaver.Implements[SecurityService]
	//Mongo
	db                components.NoSQLDatabase
	orderService      weaver.Ref[OrderService]
	orderOtherService weaver.Ref[OrderOtherService]
	roles             []string
}

func (ssi *securityService) FindAllSecurityConfigs(ctx context.Context, token string) ([]util.SecurityConfig, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return nil, err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("securityConfig")
	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var scs []util.SecurityConfig
	err = result.All(&scs)
	if err != nil {
		return nil, err
	}
	return scs, nil
}

func (ssi *securityService) Create(ctx context.Context, id, name, value, description, token string) (util.SecurityConfig, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return util.SecurityConfig{}, err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("securityConfig")

	query := fmt.Sprintf(`{"Name": %s }`, name)

	res, err := collection.FindOne(query)

	if err == nil {
		var oldSc util.SecurityConfig

		res.Decode(&oldSc)

		if oldSc.Id != "" {
			return util.SecurityConfig{}, errors.New("Security config with given name already exists!")
		}
	}

	newSc := util.SecurityConfig{
		Id:          uuid.New().String(),
		Name:        name,
		Value:       value,
		Description: description,
	}
	err = collection.InsertOne(newSc)

	if err != nil {
		return util.SecurityConfig{}, err
	}

	return newSc, nil
}

func (ssi *securityService) Update(ctx context.Context, id, name, value, description, token string) (util.SecurityConfig, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return util.SecurityConfig{}, err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("securityConfig")

	query := fmt.Sprintf(`{"Id": %s }`, id)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.SecurityConfig{}, err
	}
	var existingSc util.SecurityConfig

	err = res.Decode(&existingSc)
	if err != nil {
		return util.SecurityConfig{}, err
	}

	updatedSc := util.SecurityConfig{
		Id:          existingSc.Id,
		Name:        name,
		Value:       value,
		Description: description,
	}

	err = collection.ReplaceOne(query, updatedSc)
	if err != nil {
		return util.SecurityConfig{}, err
	}

	return updatedSc, nil
}

func (ssi *securityService) Delete(ctx context.Context, id, token string) (string, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return "", err
	}

	collection := ssi.db.GetDatabase("ts").GetCollection("securityConfig")

	query := fmt.Sprintf(`{"Id": %s }`, id)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "Delete successful", nil
}

func (ssi *securityService) Check(ctx context.Context, accountId, token string) (string, error) {
	err := util.Authenticate(token, ssi.roles...)
	if err != nil {
		return "", err
	}

	dateFormat := "Sat Jul 26 00:00:00 2025"
	dtNow := time.Now().Format(dateFormat)

	orderResult, err := ssi.orderService.Get().SecurityInfoCheck(ctx, dtNow, accountId, token)
	if err != nil {
		return "", err
	}

	orderOtherResult, err := ssi.orderOtherService.Get().SecurityInfoCheck(ctx, dtNow, accountId, token)
	if err != nil {
		return "", err
	}

	orderInOneHour := orderResult["OrderNumInLastHour"] + orderOtherResult["OrderNumInLastHour"]
	totalValidOrders := orderResult["OrderNumOfValidOrder"] + orderOtherResult["OrderNumOfValidOrder"]

	collection := ssi.db.GetDatabase("ts").GetCollection("securityConfig")
	query := fmt.Sprintf(`{"Name": %s }`, "max_order_1_hour")

	res, _ := collection.FindOne(query)
	var maxInHourConfig util.SecurityConfig
	res.Decode(&maxInHourConfig)

	query = fmt.Sprintf(`{"Name": %s }`, "max_order_not_use")

	res, _ = collection.FindOne(query)
	var maxNotUseConfig util.SecurityConfig
	res.Decode(&maxNotUseConfig)

	oneHourLine, _ := strconv.ParseUint(maxInHourConfig.Value, 10, 32)
	totalValidLine, _ := strconv.ParseFloat(maxNotUseConfig.Value, 32)

	if orderInOneHour > uint16(oneHourLine) || totalValidOrders > uint16(totalValidLine) {
		return "", errors.New("Too many orders in one hour or too many valid orders in total.")
	}

	return "Sucess", nil
}
