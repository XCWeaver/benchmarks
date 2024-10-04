package services

import (
	"context"
	"errors"
	"sync"

	"trainticket/pkg/model"

	"github.com/XCWeaver/xcweaver"
)

type OrderService interface {
	/*GetTicketListByDateAndTripId(ctx context.Context, travelDate, trainNumber, token string) ([]util.Ticket, error)
	CreateNewOrder(ctx context.Context, order util.Order, token string) (util.Order, error)*/
	//AddCreateNewOrder(ctx context.Context, order model.Order, token string) (model.Order, error)
	/*QueryOrders(ctx context.Context, orderInfo util.OrderInfo, accountId, token string) ([]util.Order, error)
	QueryOrdersForRefresh(ctx context.Context, orderInfo util.OrderInfo, accountId, token string) ([]util.Order, error)
	CalculateSoldTicket(ctx context.Context, travelDate, trainNumber, token string) (util.SoldTicket, error)
	GetOrderPrice(ctx context.Context, orderId, token string) (float32, error)
	PayOrder(ctx context.Context, orderId, token string) (util.Order, error)*/
	GetOrderById(ctx context.Context, orderId, token string) (model.Order, error)
	ModifyOrder(ctx context.Context, orderId string, status uint16, token string) (model.Order, error)
	/*SecurityInfoCheck(ctx context.Context, checkDate, accountId, token string) (map[string]uint16, error)
	SaveOrderInfo(ctx context.Context, order util.Order, token string) (util.Order, error)
	UpdateOrder(ctx context.Context, order util.Order, token string) (util.Order, error)
	DeleteOrder(ctx context.Context, orderId, token string) (string, error)
	FindAllOrder(ctx context.Context, token string) ([]util.Order, error)*/
}

/*type orderServiceOptions struct {
	MongoAddr string `toml:"mongodb_address"`
	MongoPort string `toml:"mongodb_port"`
}*/

type orderService struct {
	xcweaver.Implements[OrderService]
	//xcweaver.WithConfig[orderServiceOptions]
	//client *mongo.Client
	//stationService xcweaver.Ref[StationService]
	mu     sync.Mutex
	roles  []string
	orders []model.Order
}

func (osi *orderService) Init(ctx context.Context) error {
	logger := osi.Logger(ctx)

	/*var err error
	clientOptions := options.Client().ApplyURI("mongodb://" + osi.Config().MongoAddr + ":" + osi.Config().MongoPort + "/?directConnection=true")
	osi.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}*/

	osi.roles = append(osi.roles, "role1")
	osi.roles = append(osi.roles, "role2")
	osi.roles = append(osi.roles, "role3")

	osi.orders = initOrders()

	logger.Info("order service running!", "firstOrder", osi.orders[0])
	//logger.Info("order service running!", "mongodb_addr", osi.Config().MongoAddr, "mongodb_port", osi.Config().MongoPort)
	return nil
}

/*func (osi *orderService) GetTicketListByDateAndTripId(ctx context.Context, travelDate, trainNumber, token string) ([]util.Ticket, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"TravelDate": %s, "TrainNumber": %s}`, travelDate, trainNumber)
	res, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var orders []util.Order
	var tickets []util.Ticket

	err = res.All(&orders)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		tickets = append(tickets, util.Ticket{
			SeatNo:       order.SeatNumber,
			StartStation: order.From,
			DestStation:  order.To,
		})
	}

	return tickets, nil
}

func (osi *orderService) CreateNewOrder(ctx context.Context, order util.Order, token string) (util.Order, error) {
	err := util.Authenticate(token, osi.roles...)
	if err != nil {
		return util.Order{}, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"AccountId": %s}`, order.AccountId)
	res, err := collection.FindOne(query)
	if err == nil {
		var exOrder util.Order
		res.Decode(&exOrder)
		if exOrder.Id != "" {
			return util.Order{}, errors.New("util.Order already exists for this account.")
		}
	}

	order.Id = uuid.New().String()
	err = collection.InsertOne(order)
	if err != nil {
		return util.Order{}, nil
	}

	return order, nil
}*/

/*func (osi *orderService) AddCreateNewOrder(ctx context.Context, order model.Order, token string) (model.Order, error) {
	logger := osi.Logger(ctx)
	logger.Info("entering AddCreateNewOrder", "orderId", order.Id)

	err := util.Authenticate(token, osi.roles[0])
	if err != nil {
		return model.Order{}, err
	}

	collection := osi.client.Database("ts").Collection("orders")
	filter := bson.D{{"accountId", order.AccountId}, {"boughtDate", order.BoughtDate}, {"travelDate", order.TravelDate}, {"contactsName", order.ContactsName}, {"documentType", order.DocumentType},
		{"contactsDocumentNumber", order.ContactsDocumentNumber}, {"trainNumber", order.TrainNumber}, {"coachNumber", order.CoachNumber}, {"seatClass", order.SeatClass}, {"seatNumber", order.SeatNumber},
		{"from", order.From}, {"to", order.To}, {"status", order.Status}, {"price", order.Price}}

	res := collection.FindOne(ctx, filter)
	if res.Err() == nil {
		return model.Order{}, errors.New("Order already exists!")
	} else if res.Err() != mongo.ErrNoDocuments && res.Err() != nil {
		return model.Order{}, res.Err()
	}

	order.Id = uuid.New().String()
	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		return model.Order{}, err
	}
	logger.Debug("inserted order", "objectid", result.InsertedID)
	logger.Info("order successfully created!", "orderId", order.Id, "BoughtDate", order.BoughtDate, "TravelDate", order.TravelDate, "ContactsName", order.ContactsName,
		"DocumentType", order.DocumentType, "ContactsDocumentNumber", order.ContactsDocumentNumber, "TrainNumber", order.TrainNumber, "CoachNumber", order.CoachNumber,
		"SeatClass", model.SeatClass(order.SeatClass), "SeatNumber", order.SeatNumber, "from", order.From, "to", order.To, "status", model.OrderStatus(order.Status), "price", order.Price)

	return order, nil
}*/

/*func (osi *orderService) QueryOrders(ctx context.Context, orderInfo util.OrderInfo, accountId, token string) ([]util.Order, error) {
	err := util.Authenticate(token, osi.roles...)
	if err != nil {
		return nil, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"AccountId": %s}`, accountId)
	res, err := collection.FindMany(query)
	if err != nil {
		return nil, err
	}

	var orderList []util.Order
	err = res.All(&orderList)
	if err != nil {
		return nil, err
	}

	var finalList []util.Order

	if orderInfo.EnableTravelDateQuery || orderInfo.EnableBoughtDateQuery || orderInfo.EnableStateQuery {

		statePassFlag := false
		travelDatePassFlag := false
		boughtDatePassFlag := false

		for _, order := range orderList {

			if orderInfo.EnableStateQuery {
				if order.Status == orderInfo.State {
					statePassFlag = true
				}
			}

			if orderInfo.EnableTravelDateQuery {
				t1, _ := time.Parse(time.ANSIC, order.TravelDate)
				t2, _ := time.Parse(time.ANSIC, orderInfo.TravelDateEnd)
				t3, _ := time.Parse(time.ANSIC, order.TravelDate)
				t4, _ := time.Parse(time.ANSIC, orderInfo.TravelDateStart)

				if t1.Before(t2) && t3.Before(t4) {
					travelDatePassFlag = true
				}
			}

			if orderInfo.EnableBoughtDateQuery {
				t1, _ := time.Parse(time.ANSIC, order.BoughtDate)
				t2, _ := time.Parse(time.ANSIC, orderInfo.BoughtDateEnd)
				t3, _ := time.Parse(time.ANSIC, order.BoughtDate)
				t4, _ := time.Parse(time.ANSIC, orderInfo.BoughtDateStart)

				if t1.Before(t2) && t3.Before(t4) {
					travelDatePassFlag = true
				}
			}

			if statePassFlag && travelDatePassFlag && boughtDatePassFlag {
				finalList = append(finalList, order)
			}
		}
	} else {
		for _, order := range orderList {
			finalList = append(finalList, order)
		}
	}

	return finalList, nil
}

func (osi *orderService) QueryOrdersForRefresh(ctx context.Context, orderInfo util.OrderInfo, accountId, token string) ([]util.Order, error) {
	err := util.Authenticate(token, osi.roles...)
	if err != nil {
		return nil, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"AccountId": %s}`, accountId)
	res, err := collection.FindMany(query)
	if err != nil {
		return []util.Order{}, err
	}

	var orders []util.Order
	err = res.All(&orders)
	if err != nil {
		return orders, err
	}

	var stationIds []string
	for _, order := range orders {
		stationIds = append(stationIds, order.From)
		stationIds = append(stationIds, order.To)
	}

	names, err := osi.stationService.Get().QueryForNameBatch(ctx, stationIds, token)
	if err != nil {
		return orders, err
	}

	for idx, _ := range names {
		orders[idx].From = names[idx*2]
		orders[idx].To = names[idx*2+1]
	}

	return orders, nil
}

func (osi *orderService) CalculateSoldTicket(ctx context.Context, travelDate, trainNumber, token string) (util.SoldTicket, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.SoldTicket{}, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"TravelDate": %s, "TrainNumber": %s}`, travelDate, trainNumber)
	res, err := collection.FindMany(query)
	if err != nil {
		return util.SoldTicket{}, err
	}

	var orders []util.Order
	err = res.All(&orders)
	if err != nil {
		return util.SoldTicket{}, err
	}

	soldTicket := util.SoldTicket{}

	for _, order := range orders {
		if order.Status == uint16(util.Change) {
			continue
		}

		switch util.SeatClass(order.SeatClass) {
		case util.None:
			soldTicket.NoSeat += 1
		case util.Business:
			soldTicket.BusinessSeat += 1
		case util.FirstClass:
			soldTicket.FirstClassSeat += 1
		case util.SecondClass:
			soldTicket.SecondClassSeat += 1
		case util.HardSeat:
			soldTicket.HardSeat += 1
		case util.SoftSeat:
			soldTicket.SoftSeat += 1
		case util.HardBed:
			soldTicket.HardBed += 1
		case util.SoftBed:
			soldTicket.SoftBed += 1
		case util.HighSoftBed:
			soldTicket.HighSoftBed += 1

		default:
			continue
		}
	}

	return soldTicket, nil
}

func (osi *orderService) GetOrderPrice(ctx context.Context, orderId, token string) (float32, error) {
	err := util.Authenticate(token)
	if err != nil {
		return 0.0, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"Id": %s}`, orderId)
	res, err := collection.FindOne(query)
	if err != nil {
		return 0.0, err
	}

	var order util.Order
	err = res.Decode(&order)
	if err != nil {
		return 0.0, err
	}

	return order.Price, nil
}

func (osi *orderService) PayOrder(ctx context.Context, orderId, token string) (util.Order, error) {
	err := util.Authenticate(token)
	if err != nil {
		return util.Order{}, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"Id": %s}`, orderId)
	res, err := collection.FindOne(query)
	if err != nil {
		return util.Order{}, err
	}
	var order util.Order
	res.Decode(&order)
	update := fmt.Sprintf(`{"$set": {"Status": %d}}`, util.Paid)
	err = collection.UpdateOne(query, update)
	if err != nil {
		return util.Order{}, err
	}

	order.Status = uint16(util.Paid)
	return order, nil
}*/

func (osi *orderService) GetOrderById(ctx context.Context, orderId, token string) (model.Order, error) {
	logger := osi.Logger(ctx)
	logger.Info("entering GetOrderById", "orderId", orderId)
	osi.mu.Lock()
	defer osi.mu.Unlock()

	/*err := util.Authenticate(token)
	if err != nil {
		return model.Order{}, err
	}

	collection := osi.client.Database("ts").Collection("orders")
	filter := bson.D{{"id", orderId}}

	var order model.Order
	err := collection.FindOne(context.Background(), filter).Decode(&order)
	if err != nil {
		return model.Order{}, err
	}*/

	var order model.Order
	found := false
	for _, v := range osi.orders {
		if v.Id == orderId {
			found = true
			order = v
			break
		}
	}
	if !found {
		return model.Order{}, errors.New("order not found!")
	}

	logger.Info("order successfully found!", "orderId", order.Id, "BoughtDate", order.BoughtDate, "TravelDate", order.TravelDate, "ContactsName", order.ContactsName,
		"DocumentType", order.DocumentType, "ContactsDocumentNumber", order.ContactsDocumentNumber, "TrainNumber", order.TrainNumber, "CoachNumber", order.CoachNumber,
		"SeatClass", model.SeatClass(order.SeatClass), "SeatNumber", order.SeatNumber, "from", order.From, "to", order.To, "status", model.OrderStatus(order.Status), "price", order.Price)

	return order, nil
}

func (osi *orderService) ModifyOrder(ctx context.Context, orderId string, status uint16, token string) (model.Order, error) {
	logger := osi.Logger(ctx)
	logger.Info("entering ModifyOrder", "orderId", orderId, "status", status)
	osi.mu.Lock()
	defer osi.mu.Unlock()

	/*err := util.Authenticate(token, osi.roles...)
	if err != nil {
		return model.Order{}, err
	}

	collection := osi.client.Database("ts").Collection("orders")
	filter := bson.D{{"id", orderId}}
	var order model.Order
	result := collection.FindOne(context.Background(), filter)
	if result.Err() == mongo.ErrNoDocuments {
		return model.Order{}, errors.New(fmt.Sprintf("There is any order with the order id ", orderId))
	} else if result.Err() != nil {
		return model.Order{}, result.Err()
	}
	result.Decode(&order)

	update := bson.D{{"$set", bson.D{{"status", status}}}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.Order{}, err
	}

	order.Status = status*/

	var order model.Order
	found := false
	for i, v := range osi.orders {
		if v.Id == orderId {
			osi.orders[i].Status = status
			found = true
			break
		}
	}
	if !found {
		return model.Order{}, errors.New("order not found!")
	}

	return order, nil
}

/*func (osi *orderService) SecurityInfoCheck(ctx context.Context, checkDate, accountId, token string) (map[string]uint16, error) {
	err := util.Authenticate(token, osi.roles...)
	if err != nil {
		return nil, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"AccountId": %s}`, accountId)
	res, err := collection.FindMany(query) //TODO verify this query works
	if err != nil {
		return nil, err
	}

	var orders []util.Order
	res.All(&orders)
	countTotalValidOrder := uint16(0)
	countOrderInOneHour := uint16(0)

	dateFrom, _ := time.Parse(time.ANSIC, checkDate)

	for _, order := range orders {

		if order.Status == uint16(util.NotPaid) || order.Status == uint16(util.Paid) || order.Status == uint16(util.Collected) {
			countTotalValidOrder += 1
		}

		t1, _ := time.Parse(time.ANSIC, order.BoughtDate)

		if t1.After(dateFrom) {
			countOrderInOneHour += 1
		}
	}

	return map[string]uint16{
		"OrderNumInLastHour":   countOrderInOneHour,
		"OrderNumOfValidOrder": countTotalValidOrder,
	}, nil
}

func (osi *orderService) SaveOrderInfo(ctx context.Context, order util.Order, token string) (util.Order, error) {
	err := util.Authenticate(token, osi.roles...)
	if err != nil {
		return util.Order{}, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"Id": %s}`, order.Id)
	_, err = collection.FindOne(query)
	if err != nil {
		return util.Order{}, err
	}

	err = collection.ReplaceOne(query, order)
	if err != nil {
		return util.Order{}, nil
	}

	return order, nil
}

func (osi *orderService) UpdateOrder(ctx context.Context, order util.Order, token string) (util.Order, error) {
	err := util.Authenticate(token, osi.roles[0])
	if err != nil {
		return util.Order{}, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"Id": %s}`, order.Id)
	_, err = collection.FindOne(query)
	if err != nil {
		return util.Order{}, err
	}

	err = collection.ReplaceOne(query, order)
	if err != nil {
		return util.Order{}, nil
	}

	return order, nil
}

func (osi *orderService) DeleteOrder(ctx context.Context, orderId, token string) (string, error) {
	err := util.Authenticate(token, osi.roles...)
	if err != nil {
		return "", err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	query := fmt.Sprintf(`{"Id": %s}`, orderId)
	err = collection.DeleteOne(query)
	if err != nil {
		return "", err
	}

	return "util.Order deleted.", nil
}

func (osi *orderService) FindAllOrder(ctx context.Context, token string) ([]util.Order, error) {
	err := util.Authenticate(token)
	if err != nil {
		return nil, err
	}

	collection := osi.db.GetDatabase("ts").GetCollection("orders")

	res, err := collection.FindMany("") //TODO verify this query works
	if err != nil {
		return nil, err
	}

	var orders []util.Order
	err = res.All(&orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}*/
