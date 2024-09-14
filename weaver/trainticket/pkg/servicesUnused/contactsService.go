package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
)

type ContactService interface {
	GetAllContacts(ctx context.Context, token string) ([]model.Contact, error)
	CreateNewContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error)
	CreateNewContactsAdmin(ctx context.Context, contact model.Contact, token string) (model.Contact, error)
	DeleteContacts(ctx context.Context, contactId, token string) (string, error)
	ModifyContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error)
	FindContactsByAccountId(ctx context.Context, accountId, token string) ([]model.Contact, error)
	GetContactsByContactId(ctx context.Context, Id, token string) (model.Contact, error)
}

/*type contactServiceOptions struct {
	MongoAddr string `toml:"mongodb_address"`
	MongoPort string `toml:"mongodb_port"`
}*/

type contactService struct {
	weaver.Implements[ContactService]
	//weaver.WithConfig[contactServiceOptions]
	clientOptions *options.ClientOptions
	roles         []string
}

func (csi *contactService) Init(ctx context.Context) error {
	//logger := csi.Logger(ctx)

	//csi.clientOptions = options.Client().ApplyURI("mongodb://" + csi.Config().MongoAddr + ":" + csi.Config().MongoPort + "/?directConnection=true")

	csi.roles = append(csi.roles, "role1")
	csi.roles = append(csi.roles, "role2")
	csi.roles = append(csi.roles, "role3")

	//logger.Info("contacts service running!", "mongodb_addr", csi.Config().MongoAddr, "mongodb_port", csi.Config().MongoPort)
	return nil
}

func (csi *contactService) GetAllContacts(ctx context.Context, token string) ([]model.Contact, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering GetAllContacts")

	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, csi.clientOptions)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("contacts")
	filter := bson.D{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var contacts []model.Contact
	if err = cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}
	logger.Debug("Get all contacts executed successfully!", "contacts", contacts)

	return contacts, nil
}

func (csi *contactService) CreateNewContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering CreateNewContacts", "contactId", contact.Id)

	err := util.Authenticate(token, csi.roles[1])

	if err != nil {
		return model.Contact{}, err
	}

	client, err := mongo.Connect(ctx, csi.clientOptions)
	if err != nil {
		return model.Contact{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("contacts")
	filter := bson.D{{"accountId", contact.AccountId}, {"documentNumber", contact.DocumentNumber}, {"documentType", contact.DocumentType}}
	result := collection.FindOne(ctx, filter)
	if result.Err() == nil {
		return model.Contact{}, errors.New("Contact already exists!")
	} else if result.Err() != mongo.ErrNoDocuments && result.Err() != nil {
		return model.Contact{}, result.Err()
	}

	contact.Id = uuid.New().String()
	res, err := collection.InsertOne(ctx, contact)
	if err != nil {
		return model.Contact{}, err
	}
	logger.Debug("contact seccessfully created!", "objectid", res.InsertedID, "contactId", contact.Id, "accountId", contact.AccountId, "name", contact.Name,
		"documentType", contact.DocumentNumber, "documentNumber", contact.DocumentNumber, "phoneNumber", contact.PhoneNumber)

	return contact, nil
}

func (csi *contactService) CreateNewContactsAdmin(ctx context.Context, contact model.Contact, token string) (model.Contact, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering CreateNewContactsAdmin", "contactId", contact.Id)

	err := util.Authenticate(token, csi.roles[0])

	if err != nil {
		return model.Contact{}, err
	}

	client, err := mongo.Connect(ctx, csi.clientOptions)
	if err != nil {
		return model.Contact{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("contacts")
	filter := bson.D{{"accountId", contact.AccountId}, {"documentNumber", contact.DocumentNumber}, {"documentType", contact.DocumentType}}
	result := collection.FindOne(ctx, filter)
	if result.Err() == nil {
		return model.Contact{}, errors.New("Contact already exists!")
	} else if result.Err() != mongo.ErrNoDocuments && result.Err() != nil {
		return model.Contact{}, result.Err()
	}

	//* for ADMIN we expect the contact to be sent along in the request
	//contact.Id = uuid.New().String()
	res, err := collection.InsertOne(ctx, contact)
	if err != nil {
		return model.Contact{}, err
	}
	logger.Debug("contact seccessfully created!", "objectid", res.InsertedID, "contactId", contact.Id, "accountId", contact.AccountId, "name", contact.Name,
		"documentType", contact.DocumentNumber, "documentNumber", contact.DocumentNumber, "phoneNumber", contact.PhoneNumber)

	return contact, nil
}

func (csi *contactService) DeleteContacts(ctx context.Context, contactId, token string) (string, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering DeleteContacts", "contactId", contactId)

	err := util.Authenticate(token, csi.roles...)
	if err != nil {
		return "", err
	}

	client, err := mongo.Connect(ctx, csi.clientOptions)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("contacts")
	filter := bson.D{{"id", contactId}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return "", errors.New("Contact not found!")
	} else if result.Err() != nil {
		return "", result.Err()
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return "", err
	}

	return "Contact removed successfully", nil
}

func (csi *contactService) ModifyContacts(ctx context.Context, contact model.Contact, token string) (model.Contact, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering ModifyContacts", "contactId", contact.Id)

	err := util.Authenticate(token, csi.roles...)

	if err != nil {
		return model.Contact{}, err
	}

	client, err := mongo.Connect(ctx, csi.clientOptions)
	if err != nil {
		return model.Contact{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("contacts")
	filter := bson.D{{"id", contact.Id}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return model.Contact{}, errors.New("Contact not found!")
	} else if result.Err() != nil {
		return model.Contact{}, result.Err()
	}

	update := bson.D{{"$set", bson.D{{"accountId", contact.AccountId}, {"name", contact.Name}, {"documentType", contact.DocumentType},
		{"documentNumber", contact.DocumentNumber}, {"phoneNumber", contact.PhoneNumber}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.Contact{}, err
	}

	logger.Debug("contact seccessfully updated!", "contactId", contact.Id, "accountId", contact.AccountId, "name", contact.Name, "documentType", contact.DocumentType,
		"documentNumber", contact.DocumentNumber, "phoneNumber", contact.PhoneNumber)

	return contact, nil
}

func (csi *contactService) FindContactsByAccountId(ctx context.Context, accountId, token string) ([]model.Contact, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering FindContactsByAccountId", "accountId", accountId)

	err := util.Authenticate(token)
	if err != nil {
		return []model.Contact{}, err
	}

	client, err := mongo.Connect(ctx, csi.clientOptions)
	if err != nil {
		return []model.Contact{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("contacts")
	filter := bson.D{{"accountId", accountId}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return []model.Contact{}, err
	}

	var contacts []model.Contact
	if err = cursor.All(ctx, &contacts); err != nil {
		return nil, err
	}
	logger.Debug("FindContactsByAccountId executed successfully!", "contacts", contacts)
	return contacts, nil
}

func (csi *contactService) GetContactsByContactId(ctx context.Context, Id, token string) (model.Contact, error) {
	logger := csi.Logger(ctx)
	logger.Info("entering GetContactsByContactId", "contactId", Id)

	err := util.Authenticate(token)
	if err != nil {
		return model.Contact{}, err
	}

	client, err := mongo.Connect(ctx, csi.clientOptions)
	if err != nil {
		return model.Contact{}, err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts").Collection("contacts")
	filter := bson.D{{"id", Id}}

	result := collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return model.Contact{}, errors.New("Contact not found!")
	} else if result.Err() != nil {
		return model.Contact{}, result.Err()
	}

	var contact model.Contact
	err = result.Decode(&contact)
	if err != nil {
		return model.Contact{}, err
	}

	logger.Debug("contact seccessfully found!", "contactId", contact.Id, "accountId", contact.AccountId, "name", contact.Name, "documentType", contact.DocumentType,
		"documentNumber", contact.DocumentNumber, "phoneNumber", contact.PhoneNumber)

	return contact, nil
}
