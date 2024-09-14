package services

import (
	"context"
	"fmt"

	"trainticket/pkg/model"
	"trainticket/pkg/util"

	"github.com/TiagoMalhadas/xcweaver"
)

//* ##############################################################################################################################################
//* ##                                                                                                                                          ##
//* ##          This service normally has a queue-consumer  functionality, but we found it to be disabled in the original implementation.       ##
//* ##          Thus, we did not implement it either, for the sake of kinship between the two implementations.                                  ##
//* ##                                                                                                                                          ##
//* ##############################################################################################################################################

type NotificationService interface {
	PreserveSuccess(ctx context.Context, info model.NotificationInfo, token string) error
	OrderCreateSuccess(ctx context.Context, info model.NotificationInfo, token string) error
	OrderChangedSuccess(ctx context.Context, info model.NotificationInfo, token string) error
	OrderCancelSuccess(ctx context.Context, info model.NotificationInfo, token string) error
}

type notificationService struct {
	xcweaver.Implements[NotificationService]
	emailSender string
	roles       []string
}

func (nsi *notificationService) PreserveSuccess(ctx context.Context, info model.NotificationInfo, token string) error {

	err := util.Authenticate(token, nsi.roles...)
	if err != nil {
		return err
	}

	mail := map[string]interface{}{
		"mailFrom": nsi.emailSender,
		"mailTo":   info.Email,
		"subject":  "Reservation successful",
		"model": map[string]interface{}{
			"username":      info.Username,
			"startingPlace": info.StartingPlace,
			"endPlace":      info.EndPlace,
			"startingTime":  info.StartingTime,
			"date":          info.Date,
			"seatClass":     info.SeatClass,
			"seatNumber":    info.SeatNumber,
			"price":         info.Price,
		},
	}

	fmt.Print(mail)
	return nil
}

func (nsi *notificationService) OrderCreateSuccess(ctx context.Context, info model.NotificationInfo, token string) error {

	err := util.Authenticate(token, nsi.roles...)
	if err != nil {
		return err
	}

	mail := map[string]interface{}{
		"mailFrom": nsi.emailSender,
		"mailTo":   info.Email,
		"subject":  "Successful order creation",
		"model": map[string]interface{}{
			"username":      info.Username,
			"startingPlace": info.StartingPlace,
			"endPlace":      info.EndPlace,
			"startingTime":  info.StartingTime,
			"date":          info.Date,
			"seatClass":     info.SeatClass,
			"seatNumber":    info.SeatNumber,
			"orderNumber":   info.OrderNumber,
		},
	}

	fmt.Print(mail)
	return nil
}

func (nsi *notificationService) OrderChangedSuccess(ctx context.Context, info model.NotificationInfo, token string) error {

	err := util.Authenticate(token, nsi.roles...)
	if err != nil {
		return err
	}

	mail := map[string]interface{}{
		"mailFrom": nsi.emailSender,
		"mailTo":   info.Email,
		"subject":  "Successful order update",
		"model": map[string]interface{}{
			"username":      info.Username,
			"startingPlace": info.StartingPlace,
			"endPlace":      info.EndPlace,
			"startingTime":  info.StartingTime,
			"date":          info.Date,
			"seatClass":     info.SeatClass,
			"seatNumber":    info.SeatNumber,
			"orderNumber":   info.OrderNumber,
		},
	}

	fmt.Print(mail)
	return nil
}

func (nsi *notificationService) OrderCancelSuccess(ctx context.Context, info model.NotificationInfo, token string) error {

	err := util.Authenticate(token, nsi.roles...)
	if err != nil {
		return err
	}

	mail := map[string]interface{}{
		"mailFrom": nsi.emailSender,
		"mailTo":   info.Email,
		"subject":  "Successful order cancelation",
		"model": map[string]interface{}{
			"username": info.Username,
			"price":    info.Price,
		},
	}

	fmt.Print(mail)
	return nil
}
