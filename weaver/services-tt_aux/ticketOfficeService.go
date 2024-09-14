package services

import (
	"context"
	"fmt"
	"io/ioutil"

	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type TicketOfficeService interface {
	GetRegionList(ctx context.Context, token string) (string, error)
	GetAll(ctx context.Context, token string) ([]util.Office, error)
	GetSpecificOffices(ctx context.Context, province, city, region, token string) ([]util.Office, error)
	AddOffice(ctx context.Context, province, city, region, token string, office util.Office) (string, error)
	DeleteOffice(ctx context.Context, province, city, region, officeName, token string) (string, error)
	UpdateOffice(ctx context.Context, province, city, region, oldOfficeName, token string, office util.Office) (string, error)
}

type ticketOfficeService struct {
	weaver.Implements[TicketOfficeService]
	db components.NoSQLDatabase
}

func (tosi *ticketOfficeService) GetRegionList(ctx context.Context, token string) (string, error) {

	err := util.Authenticate(token)
	if err != nil {
		return "", err
	}

	regions, err := ioutil.ReadFile("./schema.sql")
	if err != nil {
		return "", err
	}

	return string(regions), nil
}

func (tosi *ticketOfficeService) GetAll(ctx context.Context, token string) ([]util.Office, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := tosi.db.GetDatabase("ts").GetCollection("office")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var offices []util.Office
	err = result.All(&offices)
	if err != nil {
		return nil, err
	}

	return offices, nil
}

func (tosi *ticketOfficeService) GetSpecificOffices(ctx context.Context, province, city, region, token string) ([]util.Office, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := tosi.db.GetDatabase("ts").GetCollection("office")

	query := fmt.Sprintf(`{"Province": %s, "City": %s, "Region": %s}`, province, city, region)

	result, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var offices []util.Office
	err = result.All(&offices)
	if err != nil {
		return nil, err
	}

	return offices, nil
}

func (tosi *ticketOfficeService) AddOffice(ctx context.Context, province, city, region, token string, office util.Office) (string, error) {
	err := util.Authenticate(token)
	if err != nil {
		return "", err
	}

	collection := tosi.db.GetDatabase("ts").GetCollection("office")
	query := fmt.Sprintf(`{"Province": %s, "City": %s, "Region": %s}`, province, city, region)

	update := fmt.Sprintf(`{"$push": {"Offices": %v }}`, office) //TODO check if this works @choice

	err = collection.UpdateOne(query, update)
	if err != nil {
		return "", err
	}

	return "util.Office added", nil
}

func (tosi *ticketOfficeService) DeleteOffice(ctx context.Context, province, city, region, officeName, token string) (string, error) {
	err := util.Authenticate(token)
	if err != nil {
		return "", err
	}

	collection := tosi.db.GetDatabase("ts").GetCollection("office")
	query := fmt.Sprintf(`{"Province": %s, "City": %s, "Region": %s}`, province, city, region)

	update := fmt.Sprintf(`{"$pull": {"Offices": {"OfficeName": %s}}}`, officeName) //TODO check if this works @choice

	err = collection.UpdateOne(query, update)
	if err != nil {
		return "", err
	}

	return "util.Office deleted", nil
}

func (tosi *ticketOfficeService) UpdateOffice(ctx context.Context, province, city, region, oldOfficeName, token string, office util.Office) (string, error) {
	err := util.Authenticate(token)
	if err != nil {
		return "", err
	}

	collection := tosi.db.GetDatabase("ts").GetCollection("office")
	query := fmt.Sprintf(`{"Province": %s, "City": %s, "Region": %s, "Offices.OfficeName": %s}`, province, city, region, oldOfficeName)

	update := fmt.Sprintf(`{"$set": {"Offices": [%v]} }`, office) //TODO check if this works @choice

	err = collection.UpdateOne(query, update)
	if err != nil {
		return "", err
	}

	return "util.Office updated", nil
}
