package model

import (
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/dgrijalva/jwt-go"
)

type OrderTicketInfo struct {
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
}

type FoodOrder struct {
	weaver.AutoMarshal
	Id          string
	OrderId     string
	FoodType    uint16
	StationName string
	StoreName   string
	FoodName    string
	Price       float32
}

type RebookInfo struct {
	weaver.AutoMarshal
	LoginId   string
	OrderId   string
	OldTripId string
	TripId    string
	SeatType  uint16
	Date      string
}

type RoutePlanInfo struct {
	weaver.AutoMarshal
	FromStationName string
	ToStationName   string
	TravelDate      string

	ViaStationName string
	TrainType      string
}

type Seat struct {
	weaver.AutoMarshal
	TravelDate   string
	TrainNumber  string
	StartStation string
	DestStation  string
	SeatType     uint16
}

type TravelResult struct {
	weaver.AutoMarshal
	TrainType   Train
	Route       Route
	PriceConfig PriceConfig
	Prices      map[string]float32
	Percent     float32
}

type Travel struct {
	weaver.AutoMarshal
	Trip          Trip
	StartingPlace string
	EndPlace      string
	DepartureTime string
}

type TripSummary struct {
	weaver.AutoMarshal
	Trip      Trip
	Route     Route
	TrainType Train
}

type TripDetails struct {
	weaver.AutoMarshal
	ComfortClass         uint16
	EconomyClass         uint16
	StartingStation      string
	EndStation           string
	StartingTime         string
	EndTime              string
	TripId               string
	TrainTypeId          string
	PriceForComfortClass float32
	PriceForEconomyClass float32

	StopStations                  []string
	NumberOfRestTicketFirstClass  uint16
	NumberOfRestTicketSecondClass uint16
}

type NotificationInfo struct {
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
}

type SecurityConfig struct {
	weaver.AutoMarshal
	Id          string
	Name        string
	Value       string
	Description string
}

type AddMoney struct {
	weaver.AutoMarshal
	Id     string `bson:"id"`
	UserId string `bson:"userId"`
	Money  string `bson:"money"`
	Type   string `bson:"type"`
}

type Balances struct {
	weaver.AutoMarshal
	UserId  string
	Balance float32
}

type Voucher struct {
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
}

type OrderInfo struct {
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
}

type SoldTicket struct {
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
}

type Ticket struct {
	weaver.AutoMarshal
	SeatNo       string
	StartStation string
	DestStation  string
}

type Contact struct {
	weaver.AutoMarshal
	Id             string `bson:"id"`
	AccountId      string `bson:"accountId"`
	Name           string `bson:"name"`
	DocumentType   uint16 `bson:"documentType"`
	DocumentNumber string `bson:"documentNumber"`
	PhoneNumber    string `bson:"phoneNumber"`
}

type Station struct {
	weaver.AutoMarshal
	Id       string
	Name     string
	StayTime uint16
}

type Train struct {
	weaver.AutoMarshal
	Id           string `bson:"id"`
	Name         string `bson:"name"`
	EconomyClass uint16 `bson:"economyClass"`
	ComfortClass uint16 `bson:"comfortClass"`
	AvgSpeed     uint16 `bson:"avgSpeed"`
}

type Config struct {
	weaver.AutoMarshal
	Name        string
	Value       string
	Description string
}

type PriceConfig struct {
	weaver.AutoMarshal
	Id                  string  `bson:"id"`
	TrainType           string  `bson:"trainType"`
	RouteId             string  `bson:"routeId"`
	BasicPriceRate      float32 `bson:"basicPriceRate"`
	FirstClassPriceRate float32 `bson:"firstClassPriceRate"`
}

type Consign struct {
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
}

type ConsignPriceConfig struct {
	weaver.AutoMarshal
	Id            string
	Index         uint16
	InitialWeight float32
	InitialPrice  float32
	WithinPrice   float32
	BeyondPrice   float32
}

type Order struct {
	weaver.AutoMarshal
	Id                     string  `bson:"id"`
	BoughtDate             string  `bson:"boughtDate"`
	TravelDate             string  `bson:"travelDate"`
	AccountId              string  `bson:"accountId"`
	ContactsName           string  `bson:"contactsName"`
	DocumentType           uint16  `bson:"documentType"`
	ContactsDocumentNumber string  `bson:"contactsDocumentNumber"`
	TrainNumber            string  `bson:"trainNumber"`
	CoachNumber            string  `bson:"coachNumber"`
	SeatClass              uint16  `bson:"seatClass"`
	SeatNumber             string  `bson:"seatNumber"`
	From                   string  `bson:"from"`
	To                     string  `bson:"to"`
	Status                 uint16  `bson:"status"`
	Price                  float32 `bson:"price"`
}

type RouteRequest struct {
	weaver.AutoMarshal
	Id           string
	StartStation string
	EndStation   string
	Stations     []string
	Distances    []uint16
}

type Route struct {
	weaver.AutoMarshal
	Id                string
	StartStationId    string
	TerminalStationId string
	Stations          []string
	Distances         []uint16
}

type Trip struct {
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
}

type User struct {
	weaver.AutoMarshal
	Username       string `bson:"username"`
	Password       string `bson:"password"`
	Role           string `bson:"role"`
	UserId         string `bson:"user_id"`
	Email          string `bson:"email"`
	DocumentType   uint16 `bson:"document_type"`
	DocumentNumber string `bson:"document_number"`
	Gender         uint16 `bson:"gender"`
}

type Captcha struct {
	weaver.AutoMarshal
	Name  string
	Value string
	TTL   time.Duration
}
type Insurance struct {
	weaver.AutoMarshal
	Id      string
	OrderId string
	Type    InsuranceType
}
type InsuranceType struct {
	weaver.AutoMarshal
	Id     string
	Index  uint16
	Name   string
	Price  float32
	TypeId string
}
type Payment struct {
	weaver.AutoMarshal
	Id      string `bson:"id"`
	OrderId string `bson:"orderId"`
	UserId  string `bson:"userId"`
	Price   string `bson:"price"`
	Type    string `bson:"Type"`
}
type Store struct {
	weaver.AutoMarshal
	Id           string
	StationId    string
	StoreName    string
	Telephone    string
	BusinessTime string
	DeliveryFee  string
	FoodList     []Food
}

type Delivery struct {
	weaver.AutoMarshal
	FoodName    string
	ID          string
	StationName string
	StoreName   string
}

type TrainFood struct {
	weaver.AutoMarshal
	Id       string
	TripId   string
	FoodList []Food
}

type Food struct {
	weaver.AutoMarshal
	FoodName string
	Price    float32
}

type Office struct {
	weaver.AutoMarshal
	OfficeName string
	Address    string
	WorkTime   string
	WindowNum  uint16
}

type TokenDataAux struct {
	weaver.AutoMarshal
	UserId    string
	Username  string
	Timestamp uint64
	Ttl       uint32
	Role      string
	ExpiresAt int64
}

type TokenData struct {
	UserId    string
	Username  string
	Timestamp uint64
	Ttl       uint32
	Role      string
	jwt.StandardClaims
}

//************************************** ENUMS ******************************************

type DocumentType uint16

const (
	NoneDoc DocumentType = iota
	Id_card
	Passport
	Other
)

func (d DocumentType) String() string {
	return [...]string{"NoneDoc", "Id_card", "Passport", "Other"}[d]
}

type OrderStatus uint16

const (
	NotPaid OrderStatus = iota
	Paid
	Collected
	Change
	Cancel
	Refund
	Used
)

func (o OrderStatus) String() string {
	return [...]string{"NotPaid", "Paid", "Collected", "Change", "Cancel", "Refund", "Used"}[o]
}

type SeatClass uint16

const (
	None SeatClass = iota
	Business
	FirstClass
	SecondClass
	HardSeat
	SoftSeat
	HardBed
	SoftBed
	HighSoftBed
)

func (s SeatClass) String() string {
	return [...]string{"None", "Business", "FirstClass", "SecondClass", "HardSeat", "SoftSeat", "HardBed", "SoftBed", "HighSoftBed"}[s]
}

type PaymentType uint16

const (
	NormalPayment PaymentType = iota
	Difference
	OutsidePayment
	ExternalAndDifferencePayment
)

func (s PaymentType) String() string {
	return [...]string{"Payment", "Difference", "OutsidePayment", "ExternalAndDifferencePayment"}[s]
}

type MoneyType uint16

const (
	AddMoneyType MoneyType = iota
	DrawBackMoney
)

func (s MoneyType) String() string {
	return [...]string{"AddMoney", "DrawBackMoney"}[s]
}
