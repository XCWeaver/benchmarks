package services

import (
	"context"
	"time"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v2"
)

type VerificationCodeService interface {
	Verify(ctx context.Context, receivedCode string, cookie util.Captcha) (util.Captcha, bool, error)
	Generate(ctx context.Context, cookie util.Captcha) (util.Captcha, string, error)
}

type verificationCodeService struct {
	weaver.Implements[VerificationCodeService]
	cache  ttlcache.SimpleCache
	expiry time.Duration
}

func (vcsi *verificationCodeService) Verify(ctx context.Context, receivedCode string, cookie util.Captcha) (util.Captcha, bool, error) {
	//? How is an empty struct to be checked whether empty or not? compare to util.Captcha{} ?

	//! User sends request to Verify
	//! If they have attached a Cookie, we check in the cache and compare the given code
	//! Else, we generate new cookie and return it (BUT not add it to the cache)
	//! In the latter case, we send an "Invalid" flag to the user as well.

	cookieId := ""
	var captchaCookie util.Captcha

	if cookie == (util.Captcha{}) {

		captchaCookie = util.Captcha{
			Name:  "YsbCaptcha",
			Value: cookieId,
			TTL:   vcsi.expiry,
		}
	} else {
		cookieId = cookie.Value
	}

	entry, err := vcsi.cache.Get(cookieId)
	if err != nil {
		return captchaCookie, false, err
	}

	if entry == receivedCode {

		return captchaCookie, true, nil
	}

	return captchaCookie, false, nil
}

func (vcsi *verificationCodeService) Generate(ctx context.Context, cookie util.Captcha) (util.Captcha, string, error) {

	//! Generate new captcha
	//! Add it to the TTL cache
	//! BUT, if already in cache, return what is in cache

	resCode := util.GenerateRandomString(4)

	var captchaCookie util.Captcha

	var cookieId string

	if cookie == (util.Captcha{}) {

		cookieId = uuid.New().String()

		captchaCookie = util.Captcha{
			Name:  "YsbCaptcha",
			Value: cookieId,
			TTL:   vcsi.expiry,
		}

	} else {
		cookieId = cookie.Value
	}

	vcsi.cache.SetWithTTL(cookieId, resCode, vcsi.expiry) //* can also set TTL globally

	return captchaCookie, resCode, nil
}
