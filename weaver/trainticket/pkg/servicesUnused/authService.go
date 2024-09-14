package services

import (
	"context"
	"errors"
	"log"
	"time"
	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	//GetAllUsers(ctx context.Context, token string) ([]model.User, error)
	//DeleteUserById(ctx context.Context, userId, token string) (string, error)
	CreateDefaultUser(ctx context.Context, user model.User, token string) (string, error)
	//UpdateUser(ctx context.Context, user model.User) (model.User, error)
	Login(ctx context.Context, username, password, verificationCode string, captcha model.Captcha) (string, error)
}

/*type authServiceOptions struct {
	MongoAddr string `toml:"mongodb_address"`
	MongoPort string `toml:"mongodb_port"`
}*/

type authService struct {
	weaver.Implements[AuthService]
	//Replace by code
	//verificationCodeService weaver.Ref[VerificationCodeService]
	//weaver.WithConfig[authServiceOptions]
	clientOptions *options.ClientOptions
	roles         []string
}

func (a *authService) Init(ctx context.Context) error {
	//logger := a.Logger(ctx)

	//a.clientOptions = options.Client().ApplyURI("mongodb://" + a.Config().MongoAddr + ":" + a.Config().MongoPort + "/?directConnection=true")

	a.roles = append(a.roles, "role1")
	a.roles = append(a.roles, "role2")
	a.roles = append(a.roles, "role3")

	//logger.Info("auth service running!", "mongodb_addr", a.Config().MongoAddr, "mongodb_port", a.Config().MongoPort)
	return nil
}

func (*authService) prepareSaltedHash(password string) string {

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

/*func (asi *authService) GetAllUsers(ctx context.Context, token string) ([]model.User, error) {
	err := util.Authenticate(token, asi.roles[0])

	if err != nil {
		return nil, err
	}

	collection := asi.db.GetDatabase("ts-auth-mongo").GetCollection("user")

	result, err := collection.FindMany("")

	if err != nil {
		return nil, err
	}

	var users []model.User
	err = result.All(&users)

	return users, err
}*/

/*func (asi *authService) DeleteUserById(ctx context.Context, userId, token string) (string, error) {

	err := util.Authenticate(token, asi.roles[0])

	if err != nil {
		return "", err
	}

	collection := asi.db.GetDatabase("ts-auth-mongo").GetCollection("user")

	query := fmt.Sprintf(`{"Id": %s }`, userId)

	err = collection.DeleteOne(query)

	if err != nil {
		return "", err
	}

	return "util.User removed successfully", nil
}*/

func (a *authService) CreateDefaultUser(ctx context.Context, user model.User, token string) (string, error) {
	logger := a.Logger(ctx)
	logger.Info("entering CreateDefaultUser", "username", user.Username)

	/*err := util.Authenticate(token)

	if err != nil {
		return "", err
	}*/

	if user.Username == "" || len(user.Password) < 6 {
		return "", errors.New("Invalid username or password")
	}

	client, err := mongo.Connect(ctx, a.clientOptions)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-auth-mongo").Collection("user")

	user.Password = a.prepareSaltedHash(user.Password)

	user.Role = a.roles[1]
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	logger.Debug("inserted user", "objectid", result.InsertedID)
	logger.Debug("Default user created successfully", "username", user.Username, "password", user.Password, "role", user.Role, "userId", user.UserId, "email", user.Email,
		"documentType", user.DocumentType, "documentNumber", user.DocumentNumber, "gender", user.Gender)

	return "Default user created successfully", nil
}

/*func (asi *authService) UpdateUser(ctx context.Context, username string, password string, token string) (model.User, error) {

	err := util.Authenticate(token)

	if err != nil {
		return model.User{}, err
	}

	collection := asi.db.GetDatabase("ts-auth-mongo").GetCollection("user")

	query := fmt.Sprintf(`{"Username": %s }`, username)
	result, err := collection.FindOne(query)

	if err != nil {
		return model.User{}, err
	}

	var user model.User
	result.Decode(&user)

	query = fmt.Sprintf(`{"UserId": %s }`, user.UserId)

	err = collection.DeleteOne(query)
	if err != nil {
		return model.User{}, err
	}

	newUser := model.User{
		Username: username,
		Password: asi.prepareSaltedHash(password),
		Role:     asi.roles[1],
		UserId:   user.UserId,
	}

	err = collection.InsertOne(newUser)
	if err != nil {
		return util.User{}, err
	}

	return newUser, nil
}*/

// ! TODO
func (a *authService) Login(ctx context.Context, username, password, verificationCode string, captcha model.Captcha) (string, error) {

	/*if verificationCode != "" {
		_, validCode, err := asi.verificationCodeService.Get().Verify(ctx, verificationCode, captcha)
		if err != nil {
			return "", err
		}

		if !validCode {
			return "", errors.New("Verification failed")
		}
	}*/

	logger := a.Logger(ctx)
	logger.Info("entering Login", "username", username)

	client, err := mongo.Connect(ctx, a.clientOptions)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	collection := client.Database("ts-auth-mongo").Collection("user")

	filter := bson.D{{"username", username}}

	var user model.User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Incorrect password!")
	}

	tokenData := model.TokenDataAux{
		UserId:    user.UserId,
		Username:  user.Username,
		Timestamp: uint64(time.Now().UnixMilli()),
		Ttl:       3600,
		Role:      a.roles[0], // only one role needed here
		ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
	}

	/*tik*/
	tok := util.GenerateNewToken(tokenData)

	return tok, nil
}
