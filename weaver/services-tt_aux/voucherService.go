package services

import (
	"context"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type VoucherService interface {
	Post(ctx context.Context, orderId string, typ string, token string) (util.Voucher, error)
	GetVoucher(ctx context.Context, orderId string, token string) (util.Voucher, error)
}

type voucherService struct {
	weaver.Implements[VoucherService]
	db                components.RelationalDB
	orderService      weaver.Ref[OrderService]
	orderOtherService weaver.Ref[OrderOtherService]
}

func (vsi *voucherService) Post(ctx context.Context, orderId string, typ string, token string) (util.Voucher, error) {

	err := util.Authenticate(token)
	if err != nil {
		return util.Voucher{}, err
	}

	db, err := vsi.db.Open("user", "pass", "data")
	if err != nil {
		return util.Voucher{}, err
	}
	defer db.Close()

	findQuery := "SELECT * FROM voucher where order_id = ? LIMIT 1"
	res, err := db.Query(findQuery, orderId)
	if err != nil {
		return util.Voucher{}, err
	}

	var voucher util.Voucher

	//* only get one voucher
	if res.Next() {
		res.Scan(&voucher.VoucherId, &voucher.OrderId, &voucher.TravelDate, &voucher.ContactName, &voucher.TrainNumber, &voucher.SeatClass, &voucher.SeatNumber, &voucher.StartStation, &voucher.DestStation, &voucher.Price)
		return voucher, nil
	}

	//* Insert

	var order util.Order
	if typ == "0" {
		order, err = vsi.orderService.Get().GetOrderById(ctx, orderId, token)
	} else {
		order, err = vsi.orderOtherService.Get().GetOrderById(ctx, orderId, token)
	}

	query := "INSERT INTO voucher (order_id,travelDate,contactName,trainNumber,seatClass,seatNumber,startStation,destStation,price)VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);"

	err = db.Exec(query, order.Id, order.TravelDate, order.ContactsName, order.TrainNumber, order.SeatClass, order.SeatNumber, order.From, order.To, order.Price)
	if err != nil {
		return util.Voucher{}, err
	}

	//! TODO this part is redundant
	//? Double check this

	res, err = db.Query(findQuery, orderId)
	if err != nil {
		return util.Voucher{}, err
	}

	var insertedVoucher util.Voucher

	if res.Next() {
		res.Scan(&insertedVoucher.VoucherId, &insertedVoucher.OrderId, &insertedVoucher.TravelDate, &insertedVoucher.ContactName, &insertedVoucher.TrainNumber, &insertedVoucher.SeatClass, &insertedVoucher.SeatNumber, &insertedVoucher.StartStation, &insertedVoucher.DestStation, &insertedVoucher.Price)
		return insertedVoucher, nil
	}

	return util.Voucher{}, nil
}

func (vsi *voucherService) GetVoucher(ctx context.Context, orderId, token string) (util.Voucher, error) {

	err := util.Authenticate(token)
	if err != nil {
		return util.Voucher{}, err
	}

	db, err := vsi.db.Open("user", "pass", "data")
	if err != nil {
		return util.Voucher{}, err
	}
	defer db.Close()

	findQuery := "SELECT * FROM voucher where order_id = ? LIMIT 1"

	res, err := db.Query(findQuery, orderId)
	if err != nil {
		return util.Voucher{}, err
	}

	var voucher util.Voucher

	if res.Next() {
		res.Scan(&voucher.VoucherId, &voucher.OrderId, &voucher.TravelDate, &voucher.ContactName, &voucher.TrainNumber, &voucher.SeatClass, &voucher.SeatNumber, &voucher.StartStation, &voucher.DestStation, &voucher.Price)
		return voucher, nil
	}

	return util.Voucher{}, nil
}
