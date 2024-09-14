package services

import (
	"context"
	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/XCWeaver/xcweaver"
)

type AdminBasicInfoService interface {
	GetAllContacts(ctx context.Context, token string) ([]model.Contact, error)
	DeleteContacts(ctx context.Context, contactId, token string) (string, error)
	ModifyContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error)
	AddContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error)
	/*GetAllStations(ctx context.Context, token string) ([]util.Station, error)
	DeleteStation(ctx context.Context, station util.Station, token string) (string, error)*/
	GetAllTrains(ctx context.Context, token string) ([]model.Train, error)
	DeleteTrain(ctx context.Context, trainId string, token string) (string, error)
	ModifyTrain(ctx context.Context, train model.Train, token string) (model.Train, error)
	AddTrain(ctx context.Context, train model.Train, token string) (model.Train, error)
	/*GetAllPrices(ctx context.Context, token string) ([]util.PriceConfig, error)
	DeletePrice(ctx context.Context, pc util.PriceConfig, token string) (string, error)
	ModifyPrice(ctx context.Context, pc util.PriceConfig, token string) (util.PriceConfig, error)*/
	AddPrice(ctx context.Context, pc model.PriceConfig, token string) (model.PriceConfig, error)
	//* Token-less
	/*GetAllConfigs(ctx context.Context) ([]util.Config, error)
	DeleteConfig(ctx context.Context, name string, token string) (string, error)
	ModifyConfig(ctx context.Context, config util.Config, token string) (util.Config, error)
	AddConfig(ctx context.Context, config util.Config, token string) (util.Config, error)*/
}

type adminBasicInfoService struct {
	xcweaver.Implements[AdminBasicInfoService]
	//stationService xcweaver.Ref[StationService]
	trainService xcweaver.Ref[TrainService]
	//configService  xcweaver.Ref[ConfigService]
	priceService   xcweaver.Ref[PriceService]
	contactService xcweaver.Ref[ContactService]
	roles          []string
}

func (abis *adminBasicInfoService) GetAllContacts(ctx context.Context, token string) ([]model.Contact, error) {

	err := util.Authenticate(token, abis.roles...)

	if err != nil {
		return nil, err
	}

	return abis.contactService.Get().GetAllContacts(ctx, token)
}

func (abis *adminBasicInfoService) DeleteContacts(ctx context.Context, contactId, token string) (string, error) {

	err := util.Authenticate(token, abis.roles...)

	if err != nil {
		return "", err
	}

	return abis.contactService.Get().DeleteContacts(ctx, contactId, token)
}

func (abis *adminBasicInfoService) ModifyContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error) {

	err := util.Authenticate(token, abis.roles...)

	if err != nil {
		return model.Contact{}, err
	}

	return abis.contactService.Get().ModifyContacts(ctx, contact, token)
}

func (abis *adminBasicInfoService) AddContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error) {

	err := util.Authenticate(token, abis.roles...)

	if err != nil {
		return model.Contact{}, err
	}

	return abis.contactService.Get().CreateNewContactsAdmin(ctx, contact, token)
}

/*func (abis *adminBasicInfoService) GetAllStations(ctx context.Context, token string) ([]model.Station, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return nil, err
	}

	return abis.stationService.Get().Query(ctx, token)
}

func (abis *adminBasicInfoService) DeleteStation(ctx context.Context, station model.Station, token string) (string, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return "", err
	}

	return abis.stationService.Get().Delete(ctx, station, token)
}*/

func (abis *adminBasicInfoService) GetAllTrains(ctx context.Context, token string) ([]model.Train, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return nil, err
	}

	return abis.trainService.Get().Query(ctx, token)
}

func (abis *adminBasicInfoService) DeleteTrain(ctx context.Context, trainId string, token string) (string, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return "", err
	}

	return abis.trainService.Get().Delete(ctx, trainId, token)
}

func (abis *adminBasicInfoService) ModifyTrain(ctx context.Context, train model.Train, token string) (model.Train, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return model.Train{}, err
	}

	return abis.trainService.Get().Update(ctx, train, token)
}

func (abis *adminBasicInfoService) AddTrain(ctx context.Context, train model.Train, token string) (model.Train, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return model.Train{}, err
	}

	return abis.trainService.Get().Create(ctx, train, token)
}

/*func (abis *adminBasicInfoService) GetAllPrices(ctx context.Context, token string) ([]model.PriceConfig, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return nil, err
	}

	return abis.priceService.Get().QueryAll(ctx, token)
}

func (abis *adminBasicInfoService) DeletePrice(ctx context.Context, pc model.PriceConfig, token string) (string, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return "", err
	}

	return abis.priceService.Get().Delete(ctx, pc, token)
}

func (abis *adminBasicInfoService) ModifyPrice(ctx context.Context, pc model.PriceConfig, token string) (model.PriceConfig, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return util.PriceConfig{}, err
	}

	return abis.priceService.Get().Update(ctx, pc, token)
}*/

func (abis *adminBasicInfoService) AddPrice(ctx context.Context, pc model.PriceConfig, token string) (model.PriceConfig, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return model.PriceConfig{}, err
	}

	return abis.priceService.Get().Create(ctx, pc, token)
}

//*******************************************************************************************************

/*func (abis *adminBasicInfoService) GetAllConfigs(ctx context.Context, token string) ([]model.Config, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return nil, err
	}

	return abis.configService.Get().QueryAll(ctx)
}

func (abis *adminBasicInfoService) DeleteConfig(ctx context.Context, name, token string) (string, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return "", err
	}

	return abis.configService.Get().DeleteConfig(ctx, name)
}

func (abis *adminBasicInfoService) ModifyConfig(ctx context.Context, config model.Config, token string) (model.Config, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return model.Config{}, err
	}

	return abis.configService.Get().UpdateConfig(ctx, config)
}

func (abis *adminBasicInfoService) AddConfig(ctx context.Context, config model.Config, token string) (model.Config, error) {

	err := util.Authenticate(token, abis.roles...)
	if err != nil {
		return model.Config{}, err
	}

	return abis.configService.Get().CreateConfig(ctx, config)
}*/
