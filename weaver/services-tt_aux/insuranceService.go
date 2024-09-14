package services

import (
	"context"
	"errors"
	"fmt"

	"trainticket/pkg/util"

	"github.com/google/uuid"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"

	"github.com/ServiceWeaver/weaver"
)

const TA_INSURANCE = "TRAFFIC_ACCIDENT"

type InsuranceService interface {
	GetAllInsurances(ctx context.Context, token string) ([]util.Insurance, error)
	GetAllInsuranceTypes(ctx context.Context, token string) ([]util.InsuranceType, error)
	DeleteInsurance(ctx context.Context, insuranceId, token string) (string, error)
	DeleteInsuranceByOrderId(ctx context.Context, orderId, token string) (string, error)
	ModifyInsurance(ctx context.Context, typeIndex uint16, insuranceId, orderId, token string) (util.Insurance, error)
	CreateNewInsurance(ctx context.Context, typeIndex uint16, orderId, token string) (util.Insurance, error)
	GetInsuranceById(ctx context.Context, Id, token string) (util.Insurance, error)
	FindInsuranceByOrderId(ctx context.Context, orderId, token string) (util.Insurance, error)
}

type insuranceService struct {
	weaver.Implements[InsuranceService]
	db    components.NoSQLDatabase
	roles []string
}

func (isi *insuranceService) GetAllInsurances(ctx context.Context, token string) ([]util.Insurance, error) {
	err := util.Authenticate(token, isi.roles...)
	if err != nil {
		return nil, err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurances")

	result, err := collection.FindMany("") //TODO verify this query-string works!

	if err != nil {
		return nil, err
	}

	var insurances []util.Insurance
	err = result.All(&insurances)
	if err != nil {
		return nil, err
	}

	return insurances, nil
}

func (isi *insuranceService) GetAllInsuranceTypes(ctx context.Context, token string) ([]util.InsuranceType, error) {
	err := util.Authenticate(token, isi.roles...)
	if err != nil {
		return nil, err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurance_types")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var insuranceTypes []util.InsuranceType
	err = result.All(&insuranceTypes)
	if err != nil {
		return nil, err
	}

	return insuranceTypes, nil
}

func (isi *insuranceService) DeleteInsurance(ctx context.Context, insuranceId, token string) (string, error) {
	err := util.Authenticate(token, isi.roles...)
	if err != nil {
		return "", err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurances")
	query := fmt.Sprintf(`{"Id": %s }`, insuranceId)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Insurance removed successfully", nil
}

func (isi *insuranceService) DeleteInsuranceByOrderId(ctx context.Context, orderId, token string) (string, error) {
	err := util.Authenticate(token, isi.roles...)
	if err != nil {
		return "", err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurances")
	query := fmt.Sprintf(`{"OrderId": %s }`, orderId)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Insurance removed successfully", nil
}

func (isi *insuranceService) ModifyInsurance(ctx context.Context, typeIndex uint16, insuranceId, orderId, token string) (util.Insurance, error) {
	err := util.Authenticate(token, isi.roles...)
	if err != nil {
		return util.Insurance{}, err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurances")
	query := fmt.Sprintf(`{"OrderId": %s, "Id": %s}`, orderId, insuranceId)

	res, err := collection.FindOne(query)
	if err != nil {
		return util.Insurance{}, err
	}

	var oldInsurance util.Insurance
	res.Decode(&oldInsurance)
	if oldInsurance.Id == "" {
		return util.Insurance{}, errors.New("util.Insurance does not exist!")
	}

	insTypesCollection := isi.db.GetDatabase("ts").GetCollection("insurance_types")
	typesQuery := fmt.Sprintf(`{"Index": %s }`, typeIndex)

	res, err = insTypesCollection.FindOne(typesQuery)
	if err != nil {
		return util.Insurance{}, err
	}

	var insType util.InsuranceType
	res.Decode(&insType)

	if insType.Id == "" {
		return util.Insurance{}, errors.New("util.Insurance Type does not exist!")
	}

	newInsurance := util.Insurance{
		Id:      insuranceId,
		OrderId: orderId,
		Type:    insType,
	}

	err = collection.ReplaceOne(query, newInsurance)
	if err != nil {
		return util.Insurance{}, err
	}

	return newInsurance, nil
}

func (isi *insuranceService) CreateNewInsurance(ctx context.Context, typeIndex uint16, orderId, token string) (util.Insurance, error) {
	err := util.Authenticate(token, isi.roles...)
	if err != nil {
		return util.Insurance{}, err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurances")

	query := fmt.Sprintf(`{"OrderId": %s }`, orderId)
	res, err := collection.FindOne(query)
	if err == nil {
		var oldInsurance util.Insurance

		res.Decode(&oldInsurance)

		if oldInsurance.Id != "" {
			return util.Insurance{}, errors.New("util.Insurance already exists!")
		}
	}

	insTypesCollection := isi.db.GetDatabase("ts").GetCollection("insurance_types")
	query = fmt.Sprintf(`{"Index": %s }`, typeIndex)
	res, err = insTypesCollection.FindOne(query)
	if err != nil {
		return util.Insurance{}, err
	}

	var insType util.InsuranceType
	res.Decode(&insType)

	insurance := util.Insurance{
		Id:      uuid.New().String(),
		OrderId: orderId,
		Type:    insType,
	}

	err = collection.InsertOne(insurance)
	if err != nil {
		return util.Insurance{}, err
	}

	return insurance, nil
}

func (isi *insuranceService) GetInsuranceById(ctx context.Context, Id, token string) (util.Insurance, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Insurance{}, err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurances")
	query := fmt.Sprintf(`{"Id": %s}`, Id)

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Insurance{}, err
	}

	var insurance util.Insurance
	err = result.Decode(&insurance)
	if err != nil {
		return util.Insurance{}, err
	}

	return insurance, nil
}

func (isi *insuranceService) FindInsuranceByOrderId(ctx context.Context, orderId, token string) (util.Insurance, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Insurance{}, err
	}

	collection := isi.db.GetDatabase("ts").GetCollection("insurances")
	query := fmt.Sprintf(`{"OrderId": %s}`, orderId)

	result, err := collection.FindOne(query)
	if err != nil {
		return util.Insurance{}, err
	}

	var insurance util.Insurance
	err = result.Decode(&insurance)
	if err != nil {
		return util.Insurance{}, err
	}

	return insurance, nil
}
