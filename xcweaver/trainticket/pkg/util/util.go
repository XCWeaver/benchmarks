package util

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"trainticket/pkg/model"

	"github.com/dgrijalva/jwt-go"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var secret = "secret"

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomString(length int) string {
	return StringWithCharset(length, charset)
}

func Authenticate(token string, roles ...string) error {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid Token")
		}
		return []byte(secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &model.TokenData{}, keyFunc)

	if err != nil {
		return err
	}

	payload, ok := jwtToken.Claims.(*model.TokenData)
	if !ok {
		return errors.New("Could not parse claims!")
	}

	// // * get intersection of given roles and claims' roles
	////  m := make(map[int]bool)

	//// for _, item := range roles {
	//// 	m[item] = true
	//// }

	//// for _, item := range payload.Roles {
	//// 	if _, ok := m[item]; ok {
	//// 		return nil
	//// 	}
	//// }

	if len(roles) == 0 {
		return nil
	}

	for _, item := range roles {
		if item == payload.Role {
			return nil
		}
	}

	return errors.New("Could not authenticate due to missing permissions.")
}

/*func GenerateNewToken(tokenData model.TokenData) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	tokenStr, _ := token.SignedString([]byte(secret))

	return tokenStr
}*/

func GenerateNewToken(tokenData model.TokenDataAux) string {

	jwtTokenData := model.TokenData{UserId: tokenData.UserId,
		Username:  tokenData.Username,
		Role:      tokenData.Role,
		Timestamp: tokenData.Timestamp,
		Ttl:       tokenData.Ttl,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenData.ExpiresAt,
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtTokenData)
	tokenStr, _ := token.SignedString([]byte(secret))

	return tokenStr
}

func StringToUint16(str string) (uint16, error) {
	value, err := strconv.ParseUint(str, 16, 16)
	if err != nil {
		return 0, fmt.Errorf("error parsing document type")
	}
	return uint16(value), nil
}

func StringToFloat32(str string) (float32, error) {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing document type")
	}
	return float32(value), nil
}
