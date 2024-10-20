// Code generated by "weaver generate". DO NOT EDIT.
//go:build !ignoreWeaverGen

package model

import (
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"time"
)

// weaver.InstanceOf checks.

// weaver.Router checks.

// Local stub implementations.

// Client stub implementations.

// Note that "weaver generate" will always generate the error message below.
// Everything is okay. The error message is only relevant if you see it when
// you run "go build" or "go run".
var _ codegen.LatestVersion = codegen.Version[[0][24]struct{}](`

ERROR: You generated this file with 'weaver generate' v0.24.3 (codegen
version v0.24.0). The generated code is incompatible with the version of the
github.com/ServiceWeaver/weaver module that you're using. The weaver module
version can be found in your go.mod file or by running the following command.

    go list -m github.com/ServiceWeaver/weaver

We recommend updating the weaver module and the 'weaver generate' command by
running the following.

    go get github.com/ServiceWeaver/weaver@latest
    go install github.com/ServiceWeaver/weaver/cmd/weaver@latest

Then, re-run 'weaver generate' and re-build your code. If the problem persists,
please file an issue at https://github.com/ServiceWeaver/weaver/issues.

`)

// Server stub implementations.

// Reflect stub implementations.

// AutoMarshal implementations.

var _ codegen.AutoMarshal = (*AddMoney)(nil)

type __is_AddMoney[T ~struct {
	weaver.AutoMarshal
	Id     string "bson:\"id\""
	UserId string "bson:\"userId\""
	Money  string "bson:\"money\""
	Type   string "bson:\"type\""
}] struct{}

var _ __is_AddMoney[AddMoney]

func (x *AddMoney) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("AddMoney.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.UserId)
	enc.String(x.Money)
	enc.String(x.Type)
}

func (x *AddMoney) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("AddMoney.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.UserId = dec.String()
	x.Money = dec.String()
	x.Type = dec.String()
}

var _ codegen.AutoMarshal = (*Balances)(nil)

type __is_Balances[T ~struct {
	weaver.AutoMarshal
	UserId  string
	Balance float32
}] struct{}

var _ __is_Balances[Balances]

func (x *Balances) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Balances.WeaverMarshal: nil receiver"))
	}
	enc.String(x.UserId)
	enc.Float32(x.Balance)
}

func (x *Balances) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Balances.WeaverUnmarshal: nil receiver"))
	}
	x.UserId = dec.String()
	x.Balance = dec.Float32()
}

var _ codegen.AutoMarshal = (*Captcha)(nil)

type __is_Captcha[T ~struct {
	weaver.AutoMarshal
	Name  string
	Value string
	TTL   time.Duration
}] struct{}

var _ __is_Captcha[Captcha]

func (x *Captcha) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Captcha.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Name)
	enc.String(x.Value)
	enc.Int64((int64)(x.TTL))
}

func (x *Captcha) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Captcha.WeaverUnmarshal: nil receiver"))
	}
	x.Name = dec.String()
	x.Value = dec.String()
	*(*int64)(&x.TTL) = dec.Int64()
}

var _ codegen.AutoMarshal = (*Config)(nil)

type __is_Config[T ~struct {
	weaver.AutoMarshal
	Name        string
	Value       string
	Description string
}] struct{}

var _ __is_Config[Config]

func (x *Config) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Config.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Name)
	enc.String(x.Value)
	enc.String(x.Description)
}

func (x *Config) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Config.WeaverUnmarshal: nil receiver"))
	}
	x.Name = dec.String()
	x.Value = dec.String()
	x.Description = dec.String()
}

var _ codegen.AutoMarshal = (*Consign)(nil)

type __is_Consign[T ~struct {
	weaver.AutoMarshal
	Id         string
	OrderId    string
	AccountId  string
	HandleDate string
	TargetDate string
	From       string
	To         string
	Consignee  string
	Phone      string
	Weight     float32
	Within     bool
	Price      float32
}] struct{}

var _ __is_Consign[Consign]

func (x *Consign) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Consign.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.OrderId)
	enc.String(x.AccountId)
	enc.String(x.HandleDate)
	enc.String(x.TargetDate)
	enc.String(x.From)
	enc.String(x.To)
	enc.String(x.Consignee)
	enc.String(x.Phone)
	enc.Float32(x.Weight)
	enc.Bool(x.Within)
	enc.Float32(x.Price)
}

func (x *Consign) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Consign.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.OrderId = dec.String()
	x.AccountId = dec.String()
	x.HandleDate = dec.String()
	x.TargetDate = dec.String()
	x.From = dec.String()
	x.To = dec.String()
	x.Consignee = dec.String()
	x.Phone = dec.String()
	x.Weight = dec.Float32()
	x.Within = dec.Bool()
	x.Price = dec.Float32()
}

var _ codegen.AutoMarshal = (*ConsignPriceConfig)(nil)

type __is_ConsignPriceConfig[T ~struct {
	weaver.AutoMarshal
	Id            string
	Index         uint16
	InitialWeight float32
	InitialPrice  float32
	WithinPrice   float32
	BeyondPrice   float32
}] struct{}

var _ __is_ConsignPriceConfig[ConsignPriceConfig]

func (x *ConsignPriceConfig) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("ConsignPriceConfig.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.Uint16(x.Index)
	enc.Float32(x.InitialWeight)
	enc.Float32(x.InitialPrice)
	enc.Float32(x.WithinPrice)
	enc.Float32(x.BeyondPrice)
}

func (x *ConsignPriceConfig) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("ConsignPriceConfig.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.Index = dec.Uint16()
	x.InitialWeight = dec.Float32()
	x.InitialPrice = dec.Float32()
	x.WithinPrice = dec.Float32()
	x.BeyondPrice = dec.Float32()
}

var _ codegen.AutoMarshal = (*Contact)(nil)

type __is_Contact[T ~struct {
	weaver.AutoMarshal
	Id             string "bson:\"id\""
	AccountId      string "bson:\"accountId\""
	Name           string "bson:\"name\""
	DocumentType   uint16 "bson:\"documentType\""
	DocumentNumber string "bson:\"documentNumber\""
	PhoneNumber    string "bson:\"phoneNumber\""
}] struct{}

var _ __is_Contact[Contact]

func (x *Contact) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Contact.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.AccountId)
	enc.String(x.Name)
	enc.Uint16(x.DocumentType)
	enc.String(x.DocumentNumber)
	enc.String(x.PhoneNumber)
}

func (x *Contact) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Contact.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.AccountId = dec.String()
	x.Name = dec.String()
	x.DocumentType = dec.Uint16()
	x.DocumentNumber = dec.String()
	x.PhoneNumber = dec.String()
}

var _ codegen.AutoMarshal = (*Delivery)(nil)

type __is_Delivery[T ~struct {
	weaver.AutoMarshal
	FoodName    string
	ID          string
	StationName string
	StoreName   string
}] struct{}

var _ __is_Delivery[Delivery]

func (x *Delivery) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Delivery.WeaverMarshal: nil receiver"))
	}
	enc.String(x.FoodName)
	enc.String(x.ID)
	enc.String(x.StationName)
	enc.String(x.StoreName)
}

func (x *Delivery) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Delivery.WeaverUnmarshal: nil receiver"))
	}
	x.FoodName = dec.String()
	x.ID = dec.String()
	x.StationName = dec.String()
	x.StoreName = dec.String()
}

var _ codegen.AutoMarshal = (*Food)(nil)

type __is_Food[T ~struct {
	weaver.AutoMarshal
	FoodName string
	Price    float32
}] struct{}

var _ __is_Food[Food]

func (x *Food) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Food.WeaverMarshal: nil receiver"))
	}
	enc.String(x.FoodName)
	enc.Float32(x.Price)
}

func (x *Food) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Food.WeaverUnmarshal: nil receiver"))
	}
	x.FoodName = dec.String()
	x.Price = dec.Float32()
}

var _ codegen.AutoMarshal = (*FoodOrder)(nil)

type __is_FoodOrder[T ~struct {
	weaver.AutoMarshal
	Id          string
	OrderId     string
	FoodType    uint16
	StationName string
	StoreName   string
	FoodName    string
	Price       float32
}] struct{}

var _ __is_FoodOrder[FoodOrder]

func (x *FoodOrder) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("FoodOrder.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.OrderId)
	enc.Uint16(x.FoodType)
	enc.String(x.StationName)
	enc.String(x.StoreName)
	enc.String(x.FoodName)
	enc.Float32(x.Price)
}

func (x *FoodOrder) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("FoodOrder.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.OrderId = dec.String()
	x.FoodType = dec.Uint16()
	x.StationName = dec.String()
	x.StoreName = dec.String()
	x.FoodName = dec.String()
	x.Price = dec.Float32()
}

var _ codegen.AutoMarshal = (*Insurance)(nil)

type __is_Insurance[T ~struct {
	weaver.AutoMarshal
	Id      string
	OrderId string
	Type    InsuranceType
}] struct{}

var _ __is_Insurance[Insurance]

func (x *Insurance) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Insurance.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.OrderId)
	(x.Type).WeaverMarshal(enc)
}

func (x *Insurance) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Insurance.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.OrderId = dec.String()
	(&x.Type).WeaverUnmarshal(dec)
}

var _ codegen.AutoMarshal = (*InsuranceType)(nil)

type __is_InsuranceType[T ~struct {
	weaver.AutoMarshal
	Id     string
	Index  uint16
	Name   string
	Price  float32
	TypeId string
}] struct{}

var _ __is_InsuranceType[InsuranceType]

func (x *InsuranceType) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("InsuranceType.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.Uint16(x.Index)
	enc.String(x.Name)
	enc.Float32(x.Price)
	enc.String(x.TypeId)
}

func (x *InsuranceType) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("InsuranceType.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.Index = dec.Uint16()
	x.Name = dec.String()
	x.Price = dec.Float32()
	x.TypeId = dec.String()
}

var _ codegen.AutoMarshal = (*NotificationInfo)(nil)

type __is_NotificationInfo[T ~struct {
	weaver.AutoMarshal
	SendStatus    bool
	Email         string
	OrderNumber   string
	Username      string
	StartingPlace string
	EndPlace      string
	StartingTime  string
	Date          string
	SeatClass     string
	SeatNumber    string
	Price         string
}] struct{}

var _ __is_NotificationInfo[NotificationInfo]

func (x *NotificationInfo) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("NotificationInfo.WeaverMarshal: nil receiver"))
	}
	enc.Bool(x.SendStatus)
	enc.String(x.Email)
	enc.String(x.OrderNumber)
	enc.String(x.Username)
	enc.String(x.StartingPlace)
	enc.String(x.EndPlace)
	enc.String(x.StartingTime)
	enc.String(x.Date)
	enc.String(x.SeatClass)
	enc.String(x.SeatNumber)
	enc.String(x.Price)
}

func (x *NotificationInfo) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("NotificationInfo.WeaverUnmarshal: nil receiver"))
	}
	x.SendStatus = dec.Bool()
	x.Email = dec.String()
	x.OrderNumber = dec.String()
	x.Username = dec.String()
	x.StartingPlace = dec.String()
	x.EndPlace = dec.String()
	x.StartingTime = dec.String()
	x.Date = dec.String()
	x.SeatClass = dec.String()
	x.SeatNumber = dec.String()
	x.Price = dec.String()
}

var _ codegen.AutoMarshal = (*Office)(nil)

type __is_Office[T ~struct {
	weaver.AutoMarshal
	OfficeName string
	Address    string
	WorkTime   string
	WindowNum  uint16
}] struct{}

var _ __is_Office[Office]

func (x *Office) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Office.WeaverMarshal: nil receiver"))
	}
	enc.String(x.OfficeName)
	enc.String(x.Address)
	enc.String(x.WorkTime)
	enc.Uint16(x.WindowNum)
}

func (x *Office) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Office.WeaverUnmarshal: nil receiver"))
	}
	x.OfficeName = dec.String()
	x.Address = dec.String()
	x.WorkTime = dec.String()
	x.WindowNum = dec.Uint16()
}

var _ codegen.AutoMarshal = (*Order)(nil)

type __is_Order[T ~struct {
	weaver.AutoMarshal
	Id                     string  "bson:\"id\""
	BoughtDate             string  "bson:\"boughtDate\""
	TravelDate             string  "bson:\"travelDate\""
	AccountId              string  "bson:\"accountId\""
	ContactsName           string  "bson:\"contactsName\""
	DocumentType           uint16  "bson:\"documentType\""
	ContactsDocumentNumber string  "bson:\"contactsDocumentNumber\""
	TrainNumber            string  "bson:\"trainNumber\""
	CoachNumber            string  "bson:\"coachNumber\""
	SeatClass              uint16  "bson:\"seatClass\""
	SeatNumber             string  "bson:\"seatNumber\""
	From                   string  "bson:\"from\""
	To                     string  "bson:\"to\""
	Status                 uint16  "bson:\"status\""
	Price                  float32 "bson:\"price\""
}] struct{}

var _ __is_Order[Order]

func (x *Order) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Order.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.BoughtDate)
	enc.String(x.TravelDate)
	enc.String(x.AccountId)
	enc.String(x.ContactsName)
	enc.Uint16(x.DocumentType)
	enc.String(x.ContactsDocumentNumber)
	enc.String(x.TrainNumber)
	enc.String(x.CoachNumber)
	enc.Uint16(x.SeatClass)
	enc.String(x.SeatNumber)
	enc.String(x.From)
	enc.String(x.To)
	enc.Uint16(x.Status)
	enc.Float32(x.Price)
}

func (x *Order) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Order.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.BoughtDate = dec.String()
	x.TravelDate = dec.String()
	x.AccountId = dec.String()
	x.ContactsName = dec.String()
	x.DocumentType = dec.Uint16()
	x.ContactsDocumentNumber = dec.String()
	x.TrainNumber = dec.String()
	x.CoachNumber = dec.String()
	x.SeatClass = dec.Uint16()
	x.SeatNumber = dec.String()
	x.From = dec.String()
	x.To = dec.String()
	x.Status = dec.Uint16()
	x.Price = dec.Float32()
}

var _ codegen.AutoMarshal = (*OrderInfo)(nil)

type __is_OrderInfo[T ~struct {
	weaver.AutoMarshal
	LoginId               string
	TravelDateStart       string
	TravelDateEnd         string
	BoughtDateStart       string
	BoughtDateEnd         string
	State                 uint16
	EnableTravelDateQuery bool
	EnableBoughtDateQuery bool
	EnableStateQuery      bool
}] struct{}

var _ __is_OrderInfo[OrderInfo]

func (x *OrderInfo) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("OrderInfo.WeaverMarshal: nil receiver"))
	}
	enc.String(x.LoginId)
	enc.String(x.TravelDateStart)
	enc.String(x.TravelDateEnd)
	enc.String(x.BoughtDateStart)
	enc.String(x.BoughtDateEnd)
	enc.Uint16(x.State)
	enc.Bool(x.EnableTravelDateQuery)
	enc.Bool(x.EnableBoughtDateQuery)
	enc.Bool(x.EnableStateQuery)
}

func (x *OrderInfo) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("OrderInfo.WeaverUnmarshal: nil receiver"))
	}
	x.LoginId = dec.String()
	x.TravelDateStart = dec.String()
	x.TravelDateEnd = dec.String()
	x.BoughtDateStart = dec.String()
	x.BoughtDateEnd = dec.String()
	x.State = dec.Uint16()
	x.EnableTravelDateQuery = dec.Bool()
	x.EnableBoughtDateQuery = dec.Bool()
	x.EnableStateQuery = dec.Bool()
}

var _ codegen.AutoMarshal = (*OrderTicketInfo)(nil)

type __is_OrderTicketInfo[T ~struct {
	weaver.AutoMarshal
	AccountId      string
	ContactsId     string
	TripId         string
	SeatType       uint16
	Date           string
	From           string
	To             string
	Insurance      uint16
	FoodType       uint16
	StationName    string
	StoreName      string
	FoodName       string
	FoodPrice      float32
	HandleDate     string
	ConsigneeName  string
	ConsigneePhone string
	ConsigneWeight float32
	IsWithin       bool
}] struct{}

var _ __is_OrderTicketInfo[OrderTicketInfo]

func (x *OrderTicketInfo) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("OrderTicketInfo.WeaverMarshal: nil receiver"))
	}
	enc.String(x.AccountId)
	enc.String(x.ContactsId)
	enc.String(x.TripId)
	enc.Uint16(x.SeatType)
	enc.String(x.Date)
	enc.String(x.From)
	enc.String(x.To)
	enc.Uint16(x.Insurance)
	enc.Uint16(x.FoodType)
	enc.String(x.StationName)
	enc.String(x.StoreName)
	enc.String(x.FoodName)
	enc.Float32(x.FoodPrice)
	enc.String(x.HandleDate)
	enc.String(x.ConsigneeName)
	enc.String(x.ConsigneePhone)
	enc.Float32(x.ConsigneWeight)
	enc.Bool(x.IsWithin)
}

func (x *OrderTicketInfo) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("OrderTicketInfo.WeaverUnmarshal: nil receiver"))
	}
	x.AccountId = dec.String()
	x.ContactsId = dec.String()
	x.TripId = dec.String()
	x.SeatType = dec.Uint16()
	x.Date = dec.String()
	x.From = dec.String()
	x.To = dec.String()
	x.Insurance = dec.Uint16()
	x.FoodType = dec.Uint16()
	x.StationName = dec.String()
	x.StoreName = dec.String()
	x.FoodName = dec.String()
	x.FoodPrice = dec.Float32()
	x.HandleDate = dec.String()
	x.ConsigneeName = dec.String()
	x.ConsigneePhone = dec.String()
	x.ConsigneWeight = dec.Float32()
	x.IsWithin = dec.Bool()
}

var _ codegen.AutoMarshal = (*Payment)(nil)

type __is_Payment[T ~struct {
	weaver.AutoMarshal
	Id      string "bson:\"id\""
	OrderId string "bson:\"orderId\""
	UserId  string "bson:\"userId\""
	Price   string "bson:\"price\""
	Type    string "bson:\"Type\""
}] struct{}

var _ __is_Payment[Payment]

func (x *Payment) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Payment.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.OrderId)
	enc.String(x.UserId)
	enc.String(x.Price)
	enc.String(x.Type)
}

func (x *Payment) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Payment.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.OrderId = dec.String()
	x.UserId = dec.String()
	x.Price = dec.String()
	x.Type = dec.String()
}

var _ codegen.AutoMarshal = (*PriceConfig)(nil)

type __is_PriceConfig[T ~struct {
	weaver.AutoMarshal
	Id                  string  "bson:\"id\""
	TrainType           string  "bson:\"trainType\""
	RouteId             string  "bson:\"routeId\""
	BasicPriceRate      float32 "bson:\"basicPriceRate\""
	FirstClassPriceRate float32 "bson:\"firstClassPriceRate\""
}] struct{}

var _ __is_PriceConfig[PriceConfig]

func (x *PriceConfig) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("PriceConfig.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.TrainType)
	enc.String(x.RouteId)
	enc.Float32(x.BasicPriceRate)
	enc.Float32(x.FirstClassPriceRate)
}

func (x *PriceConfig) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("PriceConfig.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.TrainType = dec.String()
	x.RouteId = dec.String()
	x.BasicPriceRate = dec.Float32()
	x.FirstClassPriceRate = dec.Float32()
}

var _ codegen.AutoMarshal = (*RebookInfo)(nil)

type __is_RebookInfo[T ~struct {
	weaver.AutoMarshal
	LoginId   string
	OrderId   string
	OldTripId string
	TripId    string
	SeatType  uint16
	Date      string
}] struct{}

var _ __is_RebookInfo[RebookInfo]

func (x *RebookInfo) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("RebookInfo.WeaverMarshal: nil receiver"))
	}
	enc.String(x.LoginId)
	enc.String(x.OrderId)
	enc.String(x.OldTripId)
	enc.String(x.TripId)
	enc.Uint16(x.SeatType)
	enc.String(x.Date)
}

func (x *RebookInfo) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("RebookInfo.WeaverUnmarshal: nil receiver"))
	}
	x.LoginId = dec.String()
	x.OrderId = dec.String()
	x.OldTripId = dec.String()
	x.TripId = dec.String()
	x.SeatType = dec.Uint16()
	x.Date = dec.String()
}

var _ codegen.AutoMarshal = (*Route)(nil)

type __is_Route[T ~struct {
	weaver.AutoMarshal
	Id                string
	StartStationId    string
	TerminalStationId string
	Stations          []string
	Distances         []uint16
}] struct{}

var _ __is_Route[Route]

func (x *Route) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Route.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.StartStationId)
	enc.String(x.TerminalStationId)
	serviceweaver_enc_slice_string_4af10117(enc, x.Stations)
	serviceweaver_enc_slice_uint16_c4441a65(enc, x.Distances)
}

func (x *Route) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Route.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.StartStationId = dec.String()
	x.TerminalStationId = dec.String()
	x.Stations = serviceweaver_dec_slice_string_4af10117(dec)
	x.Distances = serviceweaver_dec_slice_uint16_c4441a65(dec)
}

func serviceweaver_enc_slice_string_4af10117(enc *codegen.Encoder, arg []string) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		enc.String(arg[i])
	}
}

func serviceweaver_dec_slice_string_4af10117(dec *codegen.Decoder) []string {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = dec.String()
	}
	return res
}

func serviceweaver_enc_slice_uint16_c4441a65(enc *codegen.Encoder, arg []uint16) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		enc.Uint16(arg[i])
	}
}

func serviceweaver_dec_slice_uint16_c4441a65(dec *codegen.Decoder) []uint16 {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]uint16, n)
	for i := 0; i < n; i++ {
		res[i] = dec.Uint16()
	}
	return res
}

var _ codegen.AutoMarshal = (*RoutePlanInfo)(nil)

type __is_RoutePlanInfo[T ~struct {
	weaver.AutoMarshal
	FromStationName string
	ToStationName   string
	TravelDate      string
	ViaStationName  string
	TrainType       string
}] struct{}

var _ __is_RoutePlanInfo[RoutePlanInfo]

func (x *RoutePlanInfo) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("RoutePlanInfo.WeaverMarshal: nil receiver"))
	}
	enc.String(x.FromStationName)
	enc.String(x.ToStationName)
	enc.String(x.TravelDate)
	enc.String(x.ViaStationName)
	enc.String(x.TrainType)
}

func (x *RoutePlanInfo) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("RoutePlanInfo.WeaverUnmarshal: nil receiver"))
	}
	x.FromStationName = dec.String()
	x.ToStationName = dec.String()
	x.TravelDate = dec.String()
	x.ViaStationName = dec.String()
	x.TrainType = dec.String()
}

var _ codegen.AutoMarshal = (*RouteRequest)(nil)

type __is_RouteRequest[T ~struct {
	weaver.AutoMarshal
	Id           string
	StartStation string
	EndStation   string
	Stations     []string
	Distances    []uint16
}] struct{}

var _ __is_RouteRequest[RouteRequest]

func (x *RouteRequest) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("RouteRequest.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.StartStation)
	enc.String(x.EndStation)
	serviceweaver_enc_slice_string_4af10117(enc, x.Stations)
	serviceweaver_enc_slice_uint16_c4441a65(enc, x.Distances)
}

func (x *RouteRequest) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("RouteRequest.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.StartStation = dec.String()
	x.EndStation = dec.String()
	x.Stations = serviceweaver_dec_slice_string_4af10117(dec)
	x.Distances = serviceweaver_dec_slice_uint16_c4441a65(dec)
}

var _ codegen.AutoMarshal = (*Seat)(nil)

type __is_Seat[T ~struct {
	weaver.AutoMarshal
	TravelDate   string
	TrainNumber  string
	StartStation string
	DestStation  string
	SeatType     uint16
}] struct{}

var _ __is_Seat[Seat]

func (x *Seat) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Seat.WeaverMarshal: nil receiver"))
	}
	enc.String(x.TravelDate)
	enc.String(x.TrainNumber)
	enc.String(x.StartStation)
	enc.String(x.DestStation)
	enc.Uint16(x.SeatType)
}

func (x *Seat) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Seat.WeaverUnmarshal: nil receiver"))
	}
	x.TravelDate = dec.String()
	x.TrainNumber = dec.String()
	x.StartStation = dec.String()
	x.DestStation = dec.String()
	x.SeatType = dec.Uint16()
}

var _ codegen.AutoMarshal = (*SecurityConfig)(nil)

type __is_SecurityConfig[T ~struct {
	weaver.AutoMarshal
	Id          string
	Name        string
	Value       string
	Description string
}] struct{}

var _ __is_SecurityConfig[SecurityConfig]

func (x *SecurityConfig) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("SecurityConfig.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.Name)
	enc.String(x.Value)
	enc.String(x.Description)
}

func (x *SecurityConfig) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("SecurityConfig.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.Name = dec.String()
	x.Value = dec.String()
	x.Description = dec.String()
}

var _ codegen.AutoMarshal = (*SoldTicket)(nil)

type __is_SoldTicket[T ~struct {
	weaver.AutoMarshal
	TravelDate      string
	TrainNumber     string
	NoSeat          uint16
	BusinessSeat    uint16
	FirstClassSeat  uint16
	SecondClassSeat uint16
	HardSeat        uint16
	SoftSeat        uint16
	HardBed         uint16
	SoftBed         uint16
	HighSoftBed     uint16
}] struct{}

var _ __is_SoldTicket[SoldTicket]

func (x *SoldTicket) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("SoldTicket.WeaverMarshal: nil receiver"))
	}
	enc.String(x.TravelDate)
	enc.String(x.TrainNumber)
	enc.Uint16(x.NoSeat)
	enc.Uint16(x.BusinessSeat)
	enc.Uint16(x.FirstClassSeat)
	enc.Uint16(x.SecondClassSeat)
	enc.Uint16(x.HardSeat)
	enc.Uint16(x.SoftSeat)
	enc.Uint16(x.HardBed)
	enc.Uint16(x.SoftBed)
	enc.Uint16(x.HighSoftBed)
}

func (x *SoldTicket) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("SoldTicket.WeaverUnmarshal: nil receiver"))
	}
	x.TravelDate = dec.String()
	x.TrainNumber = dec.String()
	x.NoSeat = dec.Uint16()
	x.BusinessSeat = dec.Uint16()
	x.FirstClassSeat = dec.Uint16()
	x.SecondClassSeat = dec.Uint16()
	x.HardSeat = dec.Uint16()
	x.SoftSeat = dec.Uint16()
	x.HardBed = dec.Uint16()
	x.SoftBed = dec.Uint16()
	x.HighSoftBed = dec.Uint16()
}

var _ codegen.AutoMarshal = (*Station)(nil)

type __is_Station[T ~struct {
	weaver.AutoMarshal
	Id       string
	Name     string
	StayTime uint16
}] struct{}

var _ __is_Station[Station]

func (x *Station) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Station.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.Name)
	enc.Uint16(x.StayTime)
}

func (x *Station) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Station.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.Name = dec.String()
	x.StayTime = dec.Uint16()
}

var _ codegen.AutoMarshal = (*Store)(nil)

type __is_Store[T ~struct {
	weaver.AutoMarshal
	Id           string
	StationId    string
	StoreName    string
	Telephone    string
	BusinessTime string
	DeliveryFee  string
	FoodList     []Food
}] struct{}

var _ __is_Store[Store]

func (x *Store) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Store.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.StationId)
	enc.String(x.StoreName)
	enc.String(x.Telephone)
	enc.String(x.BusinessTime)
	enc.String(x.DeliveryFee)
	serviceweaver_enc_slice_Food_7a50f59f(enc, x.FoodList)
}

func (x *Store) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Store.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.StationId = dec.String()
	x.StoreName = dec.String()
	x.Telephone = dec.String()
	x.BusinessTime = dec.String()
	x.DeliveryFee = dec.String()
	x.FoodList = serviceweaver_dec_slice_Food_7a50f59f(dec)
}

func serviceweaver_enc_slice_Food_7a50f59f(enc *codegen.Encoder, arg []Food) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for i := 0; i < len(arg); i++ {
		(arg[i]).WeaverMarshal(enc)
	}
}

func serviceweaver_dec_slice_Food_7a50f59f(dec *codegen.Decoder) []Food {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make([]Food, n)
	for i := 0; i < n; i++ {
		(&res[i]).WeaverUnmarshal(dec)
	}
	return res
}

var _ codegen.AutoMarshal = (*Ticket)(nil)

type __is_Ticket[T ~struct {
	weaver.AutoMarshal
	SeatNo       string
	StartStation string
	DestStation  string
}] struct{}

var _ __is_Ticket[Ticket]

func (x *Ticket) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Ticket.WeaverMarshal: nil receiver"))
	}
	enc.String(x.SeatNo)
	enc.String(x.StartStation)
	enc.String(x.DestStation)
}

func (x *Ticket) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Ticket.WeaverUnmarshal: nil receiver"))
	}
	x.SeatNo = dec.String()
	x.StartStation = dec.String()
	x.DestStation = dec.String()
}

var _ codegen.AutoMarshal = (*TokenDataAux)(nil)

type __is_TokenDataAux[T ~struct {
	weaver.AutoMarshal
	UserId    string
	Username  string
	Timestamp uint64
	Ttl       uint32
	Role      string
	ExpiresAt int64
}] struct{}

var _ __is_TokenDataAux[TokenDataAux]

func (x *TokenDataAux) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TokenDataAux.WeaverMarshal: nil receiver"))
	}
	enc.String(x.UserId)
	enc.String(x.Username)
	enc.Uint64(x.Timestamp)
	enc.Uint32(x.Ttl)
	enc.String(x.Role)
	enc.Int64(x.ExpiresAt)
}

func (x *TokenDataAux) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TokenDataAux.WeaverUnmarshal: nil receiver"))
	}
	x.UserId = dec.String()
	x.Username = dec.String()
	x.Timestamp = dec.Uint64()
	x.Ttl = dec.Uint32()
	x.Role = dec.String()
	x.ExpiresAt = dec.Int64()
}

var _ codegen.AutoMarshal = (*Train)(nil)

type __is_Train[T ~struct {
	weaver.AutoMarshal
	Id           string "bson:\"id\""
	Name         string "bson:\"name\""
	EconomyClass uint16 "bson:\"economyClass\""
	ComfortClass uint16 "bson:\"comfortClass\""
	AvgSpeed     uint16 "bson:\"avgSpeed\""
}] struct{}

var _ __is_Train[Train]

func (x *Train) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Train.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.Name)
	enc.Uint16(x.EconomyClass)
	enc.Uint16(x.ComfortClass)
	enc.Uint16(x.AvgSpeed)
}

func (x *Train) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Train.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.Name = dec.String()
	x.EconomyClass = dec.Uint16()
	x.ComfortClass = dec.Uint16()
	x.AvgSpeed = dec.Uint16()
}

var _ codegen.AutoMarshal = (*TrainFood)(nil)

type __is_TrainFood[T ~struct {
	weaver.AutoMarshal
	Id       string
	TripId   string
	FoodList []Food
}] struct{}

var _ __is_TrainFood[TrainFood]

func (x *TrainFood) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TrainFood.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.TripId)
	serviceweaver_enc_slice_Food_7a50f59f(enc, x.FoodList)
}

func (x *TrainFood) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TrainFood.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.TripId = dec.String()
	x.FoodList = serviceweaver_dec_slice_Food_7a50f59f(dec)
}

var _ codegen.AutoMarshal = (*Travel)(nil)

type __is_Travel[T ~struct {
	weaver.AutoMarshal
	Trip          Trip
	StartingPlace string
	EndPlace      string
	DepartureTime string
}] struct{}

var _ __is_Travel[Travel]

func (x *Travel) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Travel.WeaverMarshal: nil receiver"))
	}
	(x.Trip).WeaverMarshal(enc)
	enc.String(x.StartingPlace)
	enc.String(x.EndPlace)
	enc.String(x.DepartureTime)
}

func (x *Travel) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Travel.WeaverUnmarshal: nil receiver"))
	}
	(&x.Trip).WeaverUnmarshal(dec)
	x.StartingPlace = dec.String()
	x.EndPlace = dec.String()
	x.DepartureTime = dec.String()
}

var _ codegen.AutoMarshal = (*TravelResult)(nil)

type __is_TravelResult[T ~struct {
	weaver.AutoMarshal
	TrainType   Train
	Route       Route
	PriceConfig PriceConfig
	Prices      map[string]float32
	Percent     float32
}] struct{}

var _ __is_TravelResult[TravelResult]

func (x *TravelResult) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TravelResult.WeaverMarshal: nil receiver"))
	}
	(x.TrainType).WeaverMarshal(enc)
	(x.Route).WeaverMarshal(enc)
	(x.PriceConfig).WeaverMarshal(enc)
	serviceweaver_enc_map_string_float32_359247ae(enc, x.Prices)
	enc.Float32(x.Percent)
}

func (x *TravelResult) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TravelResult.WeaverUnmarshal: nil receiver"))
	}
	(&x.TrainType).WeaverUnmarshal(dec)
	(&x.Route).WeaverUnmarshal(dec)
	(&x.PriceConfig).WeaverUnmarshal(dec)
	x.Prices = serviceweaver_dec_map_string_float32_359247ae(dec)
	x.Percent = dec.Float32()
}

func serviceweaver_enc_map_string_float32_359247ae(enc *codegen.Encoder, arg map[string]float32) {
	if arg == nil {
		enc.Len(-1)
		return
	}
	enc.Len(len(arg))
	for k, v := range arg {
		enc.String(k)
		enc.Float32(v)
	}
}

func serviceweaver_dec_map_string_float32_359247ae(dec *codegen.Decoder) map[string]float32 {
	n := dec.Len()
	if n == -1 {
		return nil
	}
	res := make(map[string]float32, n)
	var k string
	var v float32
	for i := 0; i < n; i++ {
		k = dec.String()
		v = dec.Float32()
		res[k] = v
	}
	return res
}

var _ codegen.AutoMarshal = (*Trip)(nil)

type __is_Trip[T ~struct {
	weaver.AutoMarshal
	Id                string
	TrainTypeId       string
	Number            string
	RouteId           string
	StartingTime      string
	EndTime           string
	StartingStationId string
	StationsId        string
	TerminalStationId string
}] struct{}

var _ __is_Trip[Trip]

func (x *Trip) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Trip.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Id)
	enc.String(x.TrainTypeId)
	enc.String(x.Number)
	enc.String(x.RouteId)
	enc.String(x.StartingTime)
	enc.String(x.EndTime)
	enc.String(x.StartingStationId)
	enc.String(x.StationsId)
	enc.String(x.TerminalStationId)
}

func (x *Trip) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Trip.WeaverUnmarshal: nil receiver"))
	}
	x.Id = dec.String()
	x.TrainTypeId = dec.String()
	x.Number = dec.String()
	x.RouteId = dec.String()
	x.StartingTime = dec.String()
	x.EndTime = dec.String()
	x.StartingStationId = dec.String()
	x.StationsId = dec.String()
	x.TerminalStationId = dec.String()
}

var _ codegen.AutoMarshal = (*TripDetails)(nil)

type __is_TripDetails[T ~struct {
	weaver.AutoMarshal
	ComfortClass                  uint16
	EconomyClass                  uint16
	StartingStation               string
	EndStation                    string
	StartingTime                  string
	EndTime                       string
	TripId                        string
	TrainTypeId                   string
	PriceForComfortClass          float32
	PriceForEconomyClass          float32
	StopStations                  []string
	NumberOfRestTicketFirstClass  uint16
	NumberOfRestTicketSecondClass uint16
}] struct{}

var _ __is_TripDetails[TripDetails]

func (x *TripDetails) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TripDetails.WeaverMarshal: nil receiver"))
	}
	enc.Uint16(x.ComfortClass)
	enc.Uint16(x.EconomyClass)
	enc.String(x.StartingStation)
	enc.String(x.EndStation)
	enc.String(x.StartingTime)
	enc.String(x.EndTime)
	enc.String(x.TripId)
	enc.String(x.TrainTypeId)
	enc.Float32(x.PriceForComfortClass)
	enc.Float32(x.PriceForEconomyClass)
	serviceweaver_enc_slice_string_4af10117(enc, x.StopStations)
	enc.Uint16(x.NumberOfRestTicketFirstClass)
	enc.Uint16(x.NumberOfRestTicketSecondClass)
}

func (x *TripDetails) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TripDetails.WeaverUnmarshal: nil receiver"))
	}
	x.ComfortClass = dec.Uint16()
	x.EconomyClass = dec.Uint16()
	x.StartingStation = dec.String()
	x.EndStation = dec.String()
	x.StartingTime = dec.String()
	x.EndTime = dec.String()
	x.TripId = dec.String()
	x.TrainTypeId = dec.String()
	x.PriceForComfortClass = dec.Float32()
	x.PriceForEconomyClass = dec.Float32()
	x.StopStations = serviceweaver_dec_slice_string_4af10117(dec)
	x.NumberOfRestTicketFirstClass = dec.Uint16()
	x.NumberOfRestTicketSecondClass = dec.Uint16()
}

var _ codegen.AutoMarshal = (*TripSummary)(nil)

type __is_TripSummary[T ~struct {
	weaver.AutoMarshal
	Trip      Trip
	Route     Route
	TrainType Train
}] struct{}

var _ __is_TripSummary[TripSummary]

func (x *TripSummary) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("TripSummary.WeaverMarshal: nil receiver"))
	}
	(x.Trip).WeaverMarshal(enc)
	(x.Route).WeaverMarshal(enc)
	(x.TrainType).WeaverMarshal(enc)
}

func (x *TripSummary) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("TripSummary.WeaverUnmarshal: nil receiver"))
	}
	(&x.Trip).WeaverUnmarshal(dec)
	(&x.Route).WeaverUnmarshal(dec)
	(&x.TrainType).WeaverUnmarshal(dec)
}

var _ codegen.AutoMarshal = (*User)(nil)

type __is_User[T ~struct {
	weaver.AutoMarshal
	Username       string "bson:\"username\""
	Password       string "bson:\"password\""
	Role           string "bson:\"role\""
	UserId         string "bson:\"user_id\""
	Email          string "bson:\"email\""
	DocumentType   uint16 "bson:\"document_type\""
	DocumentNumber string "bson:\"document_number\""
	Gender         uint16 "bson:\"gender\""
}] struct{}

var _ __is_User[User]

func (x *User) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("User.WeaverMarshal: nil receiver"))
	}
	enc.String(x.Username)
	enc.String(x.Password)
	enc.String(x.Role)
	enc.String(x.UserId)
	enc.String(x.Email)
	enc.Uint16(x.DocumentType)
	enc.String(x.DocumentNumber)
	enc.Uint16(x.Gender)
}

func (x *User) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("User.WeaverUnmarshal: nil receiver"))
	}
	x.Username = dec.String()
	x.Password = dec.String()
	x.Role = dec.String()
	x.UserId = dec.String()
	x.Email = dec.String()
	x.DocumentType = dec.Uint16()
	x.DocumentNumber = dec.String()
	x.Gender = dec.Uint16()
}

var _ codegen.AutoMarshal = (*Voucher)(nil)

type __is_Voucher[T ~struct {
	weaver.AutoMarshal
	VoucherId    string
	OrderId      string
	TravelDate   string
	ContactName  string
	TrainNumber  string
	SeatClass    uint16
	SeatNumber   string
	StartStation string
	DestStation  string
	Price        float32
}] struct{}

var _ __is_Voucher[Voucher]

func (x *Voucher) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("Voucher.WeaverMarshal: nil receiver"))
	}
	enc.String(x.VoucherId)
	enc.String(x.OrderId)
	enc.String(x.TravelDate)
	enc.String(x.ContactName)
	enc.String(x.TrainNumber)
	enc.Uint16(x.SeatClass)
	enc.String(x.SeatNumber)
	enc.String(x.StartStation)
	enc.String(x.DestStation)
	enc.Float32(x.Price)
}

func (x *Voucher) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("Voucher.WeaverUnmarshal: nil receiver"))
	}
	x.VoucherId = dec.String()
	x.OrderId = dec.String()
	x.TravelDate = dec.String()
	x.ContactName = dec.String()
	x.TrainNumber = dec.String()
	x.SeatClass = dec.Uint16()
	x.SeatNumber = dec.String()
	x.StartStation = dec.String()
	x.DestStation = dec.String()
	x.Price = dec.Float32()
}
