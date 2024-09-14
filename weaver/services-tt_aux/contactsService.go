package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type ContactService interface {
	GetAllContacts(ctx context.Context, token string) ([]util.Contact, error)
	CreateNewContacts(ctx context.Context, contact util.Contact, token string) (util.Contact, error)
	CreateNewContactsAdmin(ctx context.Context, contact util.Contact, token string) (util.Contact, error)
	DeleteContacts(ctx context.Context, contactId, token string) (string, error)
	ModifyContacts(ctx context.Context, contact util.Contact, token string) (util.Contact, error)
	FindContactsByAccountId(ctx context.Context, accountId, token string) (util.Contact, error)
	GetContactsByContactId(ctx context.Context, Id, token string) (util.Contact, error)
}

type contactService struct {
	weaver.Implements[ContactService]
	//MongoDB
	db    components.NoSQLDatabase
	roles []string
}

func (csi *contactService) GetAllContacts(ctx context.Context, token string) ([]util.Contact, error) {

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("contacts")

	result, err := collection.FindMany("") //TODO verify this query-string works!
	if err != nil {
		return nil, err
	}

	var contacts []util.Contact
	err = result.All(&contacts)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (csi *contactService) CreateNewContacts(ctx context.Context, contact util.Contact, token string) (util.Contact, error) {

	err := util.Authenticate(token, csi.roles[1])

	if err != nil {
		return util.Contact{}, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("contacts")

	query := fmt.Sprintf(`{"AccountId": %s }`, contact.AccountId)

	res, err := collection.FindOne(query)

	if err == nil {
		var oldContact util.Contact

		res.Decode(&oldContact)

		if oldContact.Id != "" {
			return util.Contact{}, errors.New("util.Contact already exists!")
		}
	}

	contact.Id = uuid.New().String()
	err = collection.InsertOne(contact)
	if err != nil {
		return util.Contact{}, err
	}

	return contact, nil
}

func (csi *contactService) CreateNewContactsAdmin(ctx context.Context, contact util.Contact, token string) (util.Contact, error) {

	err := util.Authenticate(token, csi.roles[0])

	if err != nil {
		return util.Contact{}, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("contacts")

	query := fmt.Sprintf(`{"Id": %s }`, contact.Id)

	res, err := collection.FindOne(query)

	if err == nil {
		var oldContact util.Contact

		res.Decode(&oldContact)

		if oldContact.Id != "" {
			return util.Contact{}, errors.New("util.Contact already exists!")
		}
	}

	//* for ADMIN we expect the contact to be sent along in the request
	//contact.Id = uuid.New().String()
	err = collection.InsertOne(contact)
	if err != nil {
		return util.Contact{}, err
	}

	return contact, nil
}

func (csi *contactService) DeleteContacts(ctx context.Context, contactId, token string) (string, error) {

	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return "", err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("contacts")
	query := fmt.Sprintf(`{"Id": %s }`, contactId)

	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Contact removed successfully", nil
}

func (csi *contactService) ModifyContacts(ctx context.Context, contact util.Contact, token string) (util.Contact, error) {

	err := util.Authenticate(token, csi.roles...)

	if err != nil {
		return util.Contact{}, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("contacts")

	query := fmt.Sprintf(`{"Id": %s }`, contact.Id)

	res, err := collection.FindOne(query)

	if err != nil {
		return util.Contact{}, err
	}

	var existingContact util.Contact
	res.Decode(&existingContact)

	if existingContact.Id == "" {
		return util.Contact{}, errors.New("Could not find contact!")
	}

	query = fmt.Sprintf(`{"Id": %s}`, contact.Id)

	err = collection.ReplaceOne(query, contact)
	if err != nil {
		return util.Contact{}, err
	}

	return contact, nil
}

func (csi *contactService) FindContactsByAccountId(ctx context.Context, accountId, token string) (util.Contact, error) {

	err := util.Authenticate(token)
	if err != nil {
		return util.Contact{}, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("contacts")

	query := fmt.Sprintf(`{"AccountId": %s}`, accountId)

	result, err := collection.FindOne(query)

	if err != nil {
		return util.Contact{}, err
	}

	var contact util.Contact
	err = result.Decode(&contact)
	if err != nil {
		return util.Contact{}, err
	}

	return contact, nil
}

func (csi *contactService) GetContactsByContactId(ctx context.Context, Id, token string) (util.Contact, error) {

	err := util.Authenticate(token)
	if err != nil {
		return util.Contact{}, err
	}

	collection := csi.db.GetDatabase("ts").GetCollection("contacts")

	query := fmt.Sprintf(`{"Id": %s}`, Id)

	result, err := collection.FindOne(query)

	if err != nil {
		return util.Contact{}, err
	}

	var contact util.Contact
	err = result.Decode(&contact)
	if err != nil {
		return util.Contact{}, err
	}

	return contact, nil
}
