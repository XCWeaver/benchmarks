package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	GetAllUsers(ctx context.Context, token string) ([]util.User, error)
	DeleteUserById(ctx context.Context, userId, token string) (string, error)
	CreateDefaultUser(ctx context.Context, user util.User, token string) (string, error)
	UpdateUser(ctx context.Context, user util.User) (util.User, error)
	Login(ctx context.Context, username, password, verificationCode string, captcha util.Captcha) (string, error)
}

type authService struct {
	weaver.Implements[AuthService]
	//Replace by code
	//Mongo
	db                      components.NoSQLDatabase
	verificationCodeService weaver.Ref[VerificationCodeService]
	roles                   []string
}

func (*authService) prepareSaltedHash(password string) string {

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

func (asi *authService) GetAllUsers(ctx context.Context, token string) ([]util.User, error) {
	err := util.Authenticate(token, asi.roles[0])

	if err != nil {
		return nil, err
	}

	collection := asi.db.GetDatabase("ts-auth-mongo").GetCollection("user")

	result, err := collection.FindMany("")

	if err != nil {
		return nil, err
	}

	var users []util.User
	err = result.All(&users)

	return users, err
}

func (asi *authService) DeleteUserById(ctx context.Context, userId, token string) (string, error) {

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
}

func (asi *authService) CreateDefaultUser(ctx context.Context, user util.User, token string) (string, error) {

	err := util.Authenticate(token)

	if err != nil {
		return "", err
	}

	if user.Username == "" || len(user.Password) < 6 {
		return "", errors.New("Invalid username or password")
	}

	collection := asi.db.GetDatabase("ts-auth-mongo").GetCollection("user")

	user.Password = asi.prepareSaltedHash(user.Password)

	user.Role = asi.roles[1]
	err = collection.InsertOne(user)
	if err != nil {
		return "", err
	}

	return "Default user created successfully", nil
}

func (asi *authService) UpdateUser(ctx context.Context, username string, password string, token string) (util.User, error) {

	err := util.Authenticate(token)

	if err != nil {
		return util.User{}, err
	}

	collection := asi.db.GetDatabase("ts-auth-mongo").GetCollection("user")

	query := fmt.Sprintf(`{"Username": %s }`, username)
	result, err := collection.FindOne(query)

	if err != nil {
		return util.User{}, err
	}

	var user util.User
	result.Decode(&user)

	query = fmt.Sprintf(`{"UserId": %s }`, user.UserId)

	err = collection.DeleteOne(query)
	if err != nil {
		return util.User{}, err
	}

	newUser := util.User{
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
}

// ! TODO
func (asi *authService) Login(ctx context.Context, username, password, verificationCode string, captcha util.Captcha) (string, error) {

	if verificationCode != "" {
		_, validCode, err := asi.verificationCodeService.Get().Verify(ctx, verificationCode, captcha)
		if err != nil {
			return "", err
		}

		if !validCode {
			return "", errors.New("Verification failed")
		}
	}

	collection := asi.db.GetDatabase("ts-auth-mongo").GetCollection("user")

	query := fmt.Sprintf(`{"Username": %s }`, username)
	result, err := collection.FindOne(query)

	if err != nil {
		return "", err
	}

	var user util.User
	result.Decode(&user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Incorrect password!")
	}

	tokenData := util.TokenData{
		UserId:    user.UserId,
		Username:  user.Username,
		Timestamp: uint64(time.Now().UnixMilli()),
		Ttl:       3600,
		Role:      asi.roles[0], // only one role needed here
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(60 * time.Minute).Unix(),
		},
	}

	/*tik*/
	tok := util.GenerateNewToken(tokenData)

	return tok, nil
}
