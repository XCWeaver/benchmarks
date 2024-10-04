package wrk2

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	tt_metrics "trainticket/pkg/metrics"
	"trainticket/pkg/services"

	"github.com/XCWeaver/xcweaver"
)

type server struct {
	xcweaver.Implements[xcweaver.Main]
	//adminBasicInfoService xcweaver.Ref[services.AdminBasicInfoService]
	//adminOrderService     xcweaver.Ref[services.AdminOrderService]
	/*adminRouteService       xcweaver.Ref[services.AdminRouteService]
	adminTravelService      xcweaver.Ref[services.AdminTravelService]*/
	//adminUserService xcweaver.Ref[services.AdminUserService]
	/*foodService             xcweaver.Ref[services.FoodService]
	routePlanService        xcweaver.Ref[services.RoutePlanService]
	consignService          xcweaver.Ref[services.ConsignService]
	voucherService          xcweaver.Ref[services.VoucherService]
	travelPlanService       xcweaver.Ref[services.TravelPlanService]
	ticketOfficeService     xcweaver.Ref[services.TicketOfficeService]
	insuranceService        xcweaver.Ref[services.InsuranceService]
	routeService            xcweaver.Ref[services.RouteService]
	rebookService xcweaver.Ref[services.RebookService]
	stationService          xcweaver.Ref[services.StationService]*/
	cancelService xcweaver.Ref[services.CancelService]
	//paymentService xcweaver.Ref[services.PaymentService]
	/*insidePaymentService    xcweaver.Ref[services.InsidePaymentService]
	notificationService     xcweaver.Ref[services.NotificationService]
	basicService            xcweaver.Ref[services.BasicService]
	ticketInfoService       xcweaver.Ref[services.TicketInfoService]
	priceService            xcweaver.Ref[services.PriceService]
	preserveOtherService    xcweaver.Ref[services.PreserveOtherService]
	preserveService         xcweaver.Ref[services.PreserveService]
	orderOtherService       xcweaver.Ref[services.OrderOtherService]
	orderService            xcweaver.Ref[services.OrderService]*/
	//contactService xcweaver.Ref[services.ContactService]
	//trainService            xcweaver.Ref[services.TrainService]
	//userService xcweaver.Ref[services.UserService]
	//authService xcweaver.Ref[services.AuthService]
	/*verificationCodeService xcweaver.Ref[services.VerificationCodeService]
	travel2Service          xcweaver.Ref[services.Travel2Service]
	travelService           xcweaver.Ref[services.TravelService]
	executeService          xcweaver.Ref[services.ExecuteService]
	securityService         xcweaver.Ref[services.SecurityService]
	configService           xcweaver.Ref[services.ConfigService]*/
	lis xcweaver.Listener `xcweaver:"wrk2"`
}

// serve is called by xcweaver.Run and contains the body of the application.
func Serve(ctx context.Context, s *server) error {

	mux := http.NewServeMux()

	// declare api endpoints
	/*mux.Handle("/wrk2-api/admin/admingGetAllContacts", instrument("admin/admingGetAllContacts", s.admingGetAllContacts, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminDeleteContacts", instrument("admin/adminDeleteContacts", s.adminDeleteContacts, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminAddContacts", instrument("admin/adminAddContacts", s.adminAddContacts, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminModifyContacts", instrument("admin/adminModifyContacts", s.adminModifyContacts, http.MethodGet, http.MethodPost))*/
	/*mux.Handle("/wrk2-api/user/adminGetAllStations", instrument("user/adminGetAllStations", s.adminGetAllStations, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/post/adminDeleteStation", instrument("post/adminDeleteStation", s.adminDeleteStation, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/home-timeline/adminGetAllTrains", instrument("home-timeline/adminGetAllTrains", s.adminGetAllTrains, http.MethodGet, http.MethodPost))*/
	/*mux.Handle("/wrk2-api/admin/adminAddTrain", instrument("admin/adminAddTrain", s.adminAddTrain, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminGetAllTrains", instrument("admin/adminGetAllTrains", s.adminGetAllTrains, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminDeleteTrain", instrument("admin/adminDeleteTrain", s.adminDeleteTrain, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminModifyTrain", instrument("admin/adminModifyTrain", s.adminModifyTrain, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminAddOrder", instrument("admin/adminAddOrder", s.adminAddOrder, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminAddPrice", instrument("admin/adminAddPrice", s.adminAddPrice, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminGetAllUsers", instrument("admin/adminGetAllUsers", s.adminGetAllUsers, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminDeleteUser", instrument("admin/adminDeleteUser", s.adminDeleteUser, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminUpdateUser", instrument("admin/adminUpdateUser", s.adminUpdateUser, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/admin/adminAddUser", instrument("admin/adminAddUser", s.adminAddUser, http.MethodGet, http.MethodPost))*/
	//mux.Handle("/wrk2-api/user-timeline/adminAddConfig", instrument("user-timeline/adminAddConfig", s.adminAddConfig, http.MethodGet, http.MethodPost))
	/*mux.Handle("/wrk2-api/user/registerUser", instrument("user/registerUser", s.registerUser, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/getAllUsers", instrument("user/getAllUsers", s.getAllUsers, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/getUserById", instrument("user/getUserById", s.getUserById, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/getUserByUsername", instrument("user/getUserByUsername", s.getUserByUsername, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/deleteUserById", instrument("user/deleteUserById", s.deleteUserById, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/updateUser", instrument("user/updateUser", s.updateUser, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/login", instrument("user/login", s.login, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/calculateRefund", instrument("user/calculateRefund", s.calculateRefund, http.MethodGet, http.MethodPost))*/
	mux.Handle("/wrk2-api/user/cancelTicket", instrument("user/cancelTicket", s.cancelTicket, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/consistencyWindow", instrument("user/consistencyWindow", s.consistencyWindow, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/inconsistencies", instrument("user/inconsistencies", s.inconsistencies, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/reset", instrument("user/reset", s.reset, http.MethodGet, http.MethodPost))
	//mux.Handle("/wrk2-api/user/pay", instrument("user/pay", s.pay, http.MethodGet, http.MethodPost))
	//mux.Handle("/wrk2-api/user/rebook", instrument("user/rebook", s.rebook, http.MethodGet, http.MethodPost))

	var handler http.Handler = mux
	s.Logger(ctx).Info("wrk2-api available", "addr", s.lis)
	return http.Serve(s.lis, handler)
}

func instrument(label string, fn func(http.ResponseWriter, *http.Request), methods ...string) http.Handler {
	allowed := map[string]struct{}{}
	for _, method := range methods {
		allowed[method] = struct{}{}
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowed[r.Method]; len(allowed) > 0 && !ok {
			msg := fmt.Sprintf("method %q not allowed", r.Method)
			http.Error(w, msg, http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	}
	return xcweaver.InstrumentHandlerFunc(label, handler)
}

/*func validateTokenParam(logger *slog.Logger, r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	// get params
	token := r.Form.Get("token")

	// validate mandatory fields
	if token == "" {
		logger.Error("you must provide a valid token")
		return "", fmt.Errorf("must provide a valid token")
	}

	return token, nil
}

func (s *server) admingGetAllContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/admingGetAllContacts")

	token, err := validateTokenParam(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	contacts, err := s.adminBasicInfoService.Get().GetAllContacts(ctx, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! list of contacts: %s\n", contacts)
}

func validateAdminDeleteContactsParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	contactID := r.Form.Get("contactID")

	// validate mandatory fields
	if token == "" {
		logger.Error("you must provide a valid token")
		return "", "", fmt.Errorf("must provide a valid token")
	}
	if contactID == "" {
		logger.Error("you must provide a contact ID")
		return "", "", fmt.Errorf("you must provide a contact ID")
	}

	return token, contactID, nil
}

func (s *server) adminDeleteContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminDeleteContacts")

	token, contactID, err := validateAdminDeleteContactsParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	removed, err := s.adminBasicInfoService.Get().DeleteContacts(ctx, contactID, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug(removed)
}

func validateAdminModifyContactsParams(logger *slog.Logger, r *http.Request) (model.Contact, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.Contact{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	contactID := r.Form.Get("contactID")
	accountID := r.Form.Get("accountID")
	name := r.Form.Get("name")
	documentType := r.Form.Get("documentType")
	documentNumber := r.Form.Get("documentNumber")
	phoneNumber := r.Form.Get("phoneNumber")

	//TO-DO
	// validate mandatory fields

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.Contact{}, "", err
	}

	return model.Contact{Id: contactID,
		AccountId:      accountID,
		Name:           name,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNumber,
		PhoneNumber:    phoneNumber}, token, nil
}

func (s *server) adminModifyContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminModifyContacts")

	contact, token, err := validateAdminModifyContactsParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	modified, err := s.adminBasicInfoService.Get().ModifyContacts(ctx, contact, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("contact %s successfully modified", modified.Id)
}

func validateAdminAddContactsParams(logger *slog.Logger, r *http.Request) (model.Contact, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.Contact{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	contactID := r.Form.Get("contactID")
	accountID := r.Form.Get("accountID")
	name := r.Form.Get("name")
	documentType := r.Form.Get("documentType")
	documentNumber := r.Form.Get("documentNumber")
	phoneNumber := r.Form.Get("phoneNumber")

	//TO-DO
	// validate mandatory fields

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.Contact{}, "", err
	}

	return model.Contact{Id: contactID,
		AccountId:      accountID,
		Name:           name,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNumber,
		PhoneNumber:    phoneNumber}, token, nil
}

func (s *server) adminAddContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminAddContacts")

	contact, token, err := validateAdminAddContactsParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	new, err := s.adminBasicInfoService.Get().AddContacts(ctx, contact, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("contact %s successfully added", new.Id)
}*/

/*func (s *server) adminGetAllStations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/adminGetAllStations")

	token, err := validateTokenParam(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	stations, err := s.adminBasicInfoService.Get().GetAllStations(ctx, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! list of stations: %s\n", stations)
}

func validateAdminDeleteStationParams(logger *slog.Logger, r *http.Request) (util.Station, string, error) {
	if err := r.ParseForm(); err != nil {
		return util.Station{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	stationID := r.Form.Get("stationID")
	stayTime := r.Form.Get("stayTime")
	name := r.Form.Get("name")

	//TO-DO
	// validate mandatory fields

	stayTimeUint16, err := util.StringToUint16(stayTime)
	if err != nil {
		logger.Error(err.Error())
		return util.Station{}, "", err
	}

	return util.Station{Id: stationID,
		Name:     name,
		StayTime: stayTimeUint16}, token, nil
}

func (s *server) adminDeleteStation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/adminDeleteStation")

	station, token, err := validateAdminDeleteStationParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	removed, err := s.adminBasicInfoService.Get().DeleteStation(ctx, station, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug(removed)
}*/

/*func (s *server) adminGetAllTrains(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/adminGetAllTrains")

	token, err := validateTokenParam(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	trains, err := s.adminBasicInfoService.Get().GetAllTrains(ctx, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! list of trains: %s\n", trains)
}

func validateAdminDeleteTrainParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	id := r.Form.Get("id")

	//TO-DO
	// validate mandatory fields

	return id, token, nil
}

func (s *server) adminDeleteTrain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminDeleteTrain")

	id, token, err := validateAdminDeleteTrainParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	result, err := s.adminBasicInfoService.Get().DeleteTrain(ctx, id, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug(result)
}

func validateAdminModifyTrainParams(logger *slog.Logger, r *http.Request) (model.Train, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.Train{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	name := r.Form.Get("name")
	economyClass := r.Form.Get("economyClass")
	comfortClass := r.Form.Get("comfortClass")
	avgSpeed := r.Form.Get("avgSpeed")

	//TO-DO
	// validate mandatory fields

	economyClassUint16, err := util.StringToUint16(economyClass)
	if err != nil {
		logger.Error(err.Error())
		return model.Train{}, "", err
	}

	comfortClassUint16, err := util.StringToUint16(comfortClass)
	if err != nil {
		logger.Error(err.Error())
		return model.Train{}, "", err
	}

	avgSpeedUint16, err := util.StringToUint16(avgSpeed)
	if err != nil {
		logger.Error(err.Error())
		return model.Train{}, "", err
	}

	return model.Train{Name: name,
		EconomyClass: economyClassUint16,
		ComfortClass: comfortClassUint16,
		AvgSpeed:     avgSpeedUint16}, token, nil
}

func (s *server) adminModifyTrain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminModifyTrain")

	train, token, err := validateAdminModifyTrainParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	train, err = s.adminBasicInfoService.Get().ModifyTrain(ctx, train, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! train modified!", "new train", train)
}

func validateAdminAddTrainParams(logger *slog.Logger, r *http.Request) (model.Train, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.Train{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	name := r.Form.Get("name")
	economyClass := r.Form.Get("economyClass")
	comfortClass := r.Form.Get("comfortClass")
	avgSpeed := r.Form.Get("avgSpeed")

	//TO-DO
	// validate mandatory fields

	economyClassUint16, err := util.StringToUint16(economyClass)
	if err != nil {
		logger.Error(err.Error())
		return model.Train{}, "", err
	}

	comfortClassUint16, err := util.StringToUint16(comfortClass)
	if err != nil {
		logger.Error(err.Error())
		return model.Train{}, "", err
	}

	avgSpeedUint16, err := util.StringToUint16(avgSpeed)
	if err != nil {
		logger.Error(err.Error())
		return model.Train{}, "", err
	}

	return model.Train{Name: name,
		EconomyClass: economyClassUint16,
		ComfortClass: comfortClassUint16,
		AvgSpeed:     avgSpeedUint16}, token, nil
}

func (s *server) adminAddTrain(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminAddTrain")

	train, token, err := validateAdminAddTrainParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	train, err = s.adminBasicInfoService.Get().AddTrain(ctx, train, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! new train created!", "train", train)
}*/

/*func (s *server) AdminGetAllConfigs(ctx context.Context, token string) ([]util.Config, error) {
	return s.adminBasicInfoService.GetAllConfigs(ctx)
}
func (s *server) AdminDeleteConfig(ctx context.Context, name string, token string) (string, error) {
	return s.adminBasicInfoService.DeleteConfig(ctx, name, token)
}
func (s *server) AdminModifyConfig(ctx context.Context, config util.Config, token string) (util.Config, error) {
	return s.adminBasicInfoService.ModifyConfig(ctx, config)
}

func validateAdminAddConfigParams(logger *slog.Logger, r *http.Request) (util.Config, string, error) {
	if err := r.ParseForm(); err != nil {
		return util.Config{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	name := r.Form.Get("name")
	value := r.Form.Get("vaue")
	description := r.Form.Get("description")

	//TO-DO
	// validate mandatory fields

	return util.Config{Name: name,
		Value:       value,
		Description: description}, token, nil
}

func (s *server) adminAddConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/adminAddConfig")

	config, token, err := validateAdminAddConfigParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	config, err = s.adminBasicInfoService.Get().AddConfig(ctx, config, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! new config created: %s\n", config)

}

func (s *server) AdminGetAllPrices(ctx context.Context, token string) ([]util.PriceConfig, error) {
	return s.adminBasicInfoService.GetAllPrices(ctx, token)
}
func (s *server) AdminDeletePrice(ctx context.Context, pc util.PriceConfig, token string) (string, error) {
	return s.adminBasicInfoService.DeletePrice(ctx, pc, token)
}
func (s *server) AdminModifyPrice(ctx context.Context, pc util.PriceConfig, token string) (util.PriceConfig, error) {
	return s.adminBasicInfoService.ModifyPrice(ctx, pc, token)
}*/

/*func validateAdminAddPriceParams(logger *slog.Logger, r *http.Request) (model.PriceConfig, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.PriceConfig{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	trainType := r.Form.Get("trainType")
	routeId := r.Form.Get("routeId")
	basicPriceRate := r.Form.Get("basicPriceRate")
	firstClassPriceRate := r.Form.Get("firstClassPriceRate")

	//TO-DO
	// validate mandatory fields

	basicPriceRateFloat32, err := util.StringToFloat32(basicPriceRate)
	if err != nil {
		logger.Error(err.Error())
		return model.PriceConfig{}, "", err
	}

	firstClassPriceRateFloat32, err := util.StringToFloat32(firstClassPriceRate)
	if err != nil {
		logger.Error(err.Error())
		return model.PriceConfig{}, "", err
	}

	return model.PriceConfig{TrainType: trainType,
		RouteId:             routeId,
		BasicPriceRate:      basicPriceRateFloat32,
		FirstClassPriceRate: firstClassPriceRateFloat32}, token, nil
}

func (s *server) adminAddPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/adminAddPrice")

	pc, token, err := validateAdminAddPriceParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	pc, err = s.adminBasicInfoService.Get().AddPrice(ctx, pc, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! new priceconfig created: %s\n", pc)
}*/

/*func (s *server) AdminGetAllOrders(ctx context.Context, token string) ([]util.Order, error) {
	return s.adminOrderService.GetAllOrders(ctx, token)
}
func (s *server) AdminDeleteOrder(ctx context.Context, orderID string, token string) (string, error) {
	return s.adminOrderService.DeleteOrder(ctx, orderID, token)
}
func (s *server) AdminUpdateOrder(ctx context.Context, order util.Order, token string) (util.Order, error) {
	return s.adminOrderService.UpdateOrder(ctx, order, token)
}*/

/*func validateAdminAddOrderParams(logger *slog.Logger, r *http.Request) (model.Order, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.Order{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	boughtDate := r.Form.Get("boughtDate")
	travelDate := r.Form.Get("travelDate")
	accountId := r.Form.Get("accountId")
	contactsName := r.Form.Get("contactsName")
	documentType := r.Form.Get("documentType")
	contactsDocumentNumber := r.Form.Get("contactsDocumentNumber")
	trainNumber := r.Form.Get("trainNumber")
	coachNumber := r.Form.Get("coachNumber")
	seatClass := r.Form.Get("seatClass")
	seatNumber := r.Form.Get("seatNumber")
	from := r.Form.Get("from")
	to := r.Form.Get("to")
	status := r.Form.Get("status")
	price := r.Form.Get("price")

	//TO-DO
	// validate mandatory fields

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.Order{}, "", err
	}

	seatClassUint16, err := util.StringToUint16(seatClass)
	if err != nil {
		logger.Error(err.Error())
		return model.Order{}, "", err
	}

	statusUint16, err := util.StringToUint16(status)
	if err != nil {
		logger.Error(err.Error())
		return model.Order{}, "", err
	}

	priceFloat32, err := util.StringToFloat32(price)
	if err != nil {
		logger.Error(err.Error())
		return model.Order{}, "", err
	}

	return model.Order{BoughtDate: boughtDate,
		TravelDate:             travelDate,
		AccountId:              accountId,
		ContactsName:           contactsName,
		DocumentType:           documentTypeUint16,
		ContactsDocumentNumber: contactsDocumentNumber,
		TrainNumber:            trainNumber,
		CoachNumber:            coachNumber,
		SeatClass:              seatClassUint16,
		SeatNumber:             seatNumber,
		From:                   from,
		To:                     to,
		Status:                 statusUint16,
		Price:                  priceFloat32}, token, nil
}

func (s *server) adminAddOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminAddOrder")
	addOrderStartMs := time.Now().UnixMilli()

	order, token, err := validateAdminAddOrderParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		w.Header().Set("Content-Type", "application/json")
		orderJson, err := json.Marshal(model.Order{})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		w.Write(orderJson)
		return
	}

	order, err = s.adminOrderService.Get().AddOrder(ctx, order, token)
	if err != nil {
		logger.Error(err.Error())
		w.Header().Set("Content-Type", "application/json")
		orderJson, err := json.Marshal(model.Order{})
		if err != nil {
			logger.Error(err.Error())
			return
		}
		w.Write(orderJson)
		return
	}

	logger.Debug("success! new order created!", "order", order)
	w.Header().Set("Content-Type", "application/json")
	orderJson, err := json.Marshal(order)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	w.Write(orderJson)
	tt_metrics.OrderTicketDuration.Put(float64(time.Now().UnixMilli() - addOrderStartMs))
	tt_metrics.Orders.Inc()
}*/

/*
	func (s *server) AdminGetAllRoutes(ctx context.Context, token string) ([]util.Route, error) {
		return s.adminRouteService.GetAllRoutes(ctx, token)
	}

	func (s *server) AdminAddRoute(ctx context.Context, id string, startStation string, endStation string, stationList []string, distanceList []uint16, token string) (util.Route, error) {
		routereq := util.RouteRequest{Id: id, StartStation: startStation, EndStation: endStation, Stations: stationList, Distances: distanceList}
		return s.adminRouteService.AddRoute(ctx, routereq, token)
	}

	func (s *server) AdminDeleteRoute(ctx context.Context, routeId string, token string) (string, error) {
		return s.adminRouteService.DeleteRoute(ctx, routeId, token)
	}

	func (s *server) AdminGetAllTravels(ctx context.Context, token string) ([]util.Trip, []util.Train, []util.Route, error) {
		return s.adminTravelService.GetAllTravels(ctx, token)
	}

	func (s *server) AdminAddTravel(ctx context.Context, tripID string, trainTypeID string, number string, routeID string, startingStationID string, stationIDs string, terminalStationId string, startingTime string, endTime string, token string) (util.Trip, error) {
		trip := util.Trip{tripID, trainTypeID, number, routeID, startingTime, endTime, startingStationID, stationIDs, terminalStationId}
		return s.adminTravelService.AddTravel(ctx, trip, token)
	}

	func (s *server) AdminUpdateTravel(ctx context.Context, tripID string, trainTypeID string, number string, routeID string, startingStationID string, stationIDs string, terminalStationId string, startingTime string, endTime string, token string) (util.Trip, error) {
		trip := util.Trip{tripID, trainTypeID, number, routeID, startingTime, endTime, startingStationID, stationIDs, terminalStationId}
		return s.adminTravelService.UpdateTravel(ctx, trip, token)
	}

	func (s *server) AdminDeleteTravel(ctx context.Context, tripID string, token string) (string, error) {
		return s.adminTravelService.DeleteTravel(ctx, tripID, token)
	}*/

/*func validateAdminGetAllUsersParams(logger *slog.Logger, r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	// get params
	token := r.Form.Get("token")

	//TO-DO
	// validate mandatory fields

	return token, nil
}

func (s *server) adminGetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminGetAllUsers")

	token, err := validateAdminGetAllUsersParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.adminUserService.Get().GetAllUsers(ctx, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

}

func validateAdminDeleteUserParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	userId := r.Form.Get("userId")

	//TO-DO
	// validate mandatory fields

	return userId, token, nil
}

func (s *server) adminDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminDeleteUser")

	userId, token, err := validateAdminDeleteUserParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.adminUserService.Get().DeleteUser(ctx, userId, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func validateAdminUpdateUserParams(logger *slog.Logger, r *http.Request) (model.User, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.User{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	gender := r.Form.Get("gender")
	documentType := r.Form.Get("documentType")
	documentNum := r.Form.Get("documentNum")
	email := r.Form.Get("email")

	//TO-DO
	// validate mandatory fields

	genderUint16, err := util.StringToUint16(gender)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	return model.User{Username: username,
		Password:       password,
		Gender:         genderUint16,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNum,
		Email:          email}, token, nil
}

func (s *server) adminUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminUpdateUser")

	user, token, err := validateAdminUpdateUserParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.adminUserService.Get().UpdateUser(ctx, user, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func validateAdminAddUserParams(logger *slog.Logger, r *http.Request) (model.User, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.User{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	gender := r.Form.Get("gender")
	documentType := r.Form.Get("documentType")
	documentNum := r.Form.Get("documentNum")
	email := r.Form.Get("email")

	//TO-DO
	// validate mandatory fields

	genderUint16, err := util.StringToUint16(gender)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	return model.User{Username: username,
		Password:       password,
		Gender:         genderUint16,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNum,
		Email:          email}, token, nil
}

func (s *server) adminAddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/admin/adminAddUser")

	user, token, err := validateAdminAddUserParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.adminUserService.Get().AddUser(ctx, user, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}*/

/*func (s *server) FindAllFoodOrder(ctx context.Context, token string) ([]util.FoodOrder, error) {
	return s.foodService.FindAllFoodOrder(ctx, token)
}

func (s *server) CreateFoodOrder(ctx context.Context, foodOrder util.FoodOrder, token string) (util.FoodOrder, error) {
	return s.foodService.CreateFoodOrder(ctx, foodOrder, token)
}

func (s *server) UpdateFoodOrder(ctx context.Context, foodOrder util.FoodOrder, token string) (util.FoodOrder, error) {
	return s.foodService.UpdateFoodOrder(ctx, foodOrder, token)
}

func (s *server) CancelFoodOrder(ctx context.Context, orderId string, token string) (string, error) {
	return s.foodService.DeleteFoodOrder(ctx, orderId, token)
}

func (s *server) FindFoodOrderByOrderId(ctx context.Context, orderId string, token string) ([]util.FoodOrder, error) {
	return s.foodService.FindFoodOrderByOrderId(ctx, orderId, token)
}

func (s *server) GetAllFoods(ctx context.Context, date string, startStation string, endStation string, tripID string, token string) ([]util.Food, map[string]util.Store, error) {
	return s.foodService.GetAllFood(ctx, startStation, endStation, tripID, token)
}

func (s *server) GetCheapestRoutes(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	return s.routePlanService.GetCheapestRoutes(ctx, info, token)
}

func (s *server) GetQuickestRoutes(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	return s.routePlanService.GetQuickestRoutes(ctx, info, token)
}

func (s *server) GetMinStopStationsRoutes(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	return s.routePlanService.GetMinStopStations(ctx, info, token)
}

func (s *server) InsertConsign(ctx context.Context, consign util.Consign, token string) (util.Consign, error) {
	return s.consignService.InsertConsign(ctx, consign, token)
}

func (s *server) UpdateConsign(ctx context.Context, consign util.Consign, token string) (util.Consign, error) {
	return s.consignService.UpdateConsign(ctx, consign, token)
}

func (s *server) FindConsignByAccountId(ctx context.Context, accountID string, token string) ([]util.Consign, error) {
	return s.consignService.FindByAccountId(ctx, accountID, token)
}

func (s *server) FindConsignByOrderId(ctx context.Context, orderID string, token string) ([]util.Consign, error) {
	return s.consignService.FindByOrderId(ctx, orderID, token)
}

func (s *server) FindConsignByConsignee(ctx context.Context, consignee string, token string) ([]util.Consign, error) {
	return s.consignService.FindByConsignee(ctx, consignee, token)
}

func (s *server) GetVoucher(ctx context.Context, orderId string, token string) (util.Voucher, error) {
	return s.voucherService.GetVoucher(ctx, orderId, token)
}

func (s *server) PostVoucher(ctx context.Context, orderId string, typ string, token string) (util.Voucher, error) {
	return s.voucherService.Post(ctx, orderId, typ, token)
}

func (s *server) GetTransferResult(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, []util.TripDetails, error) {
	return s.travelPlanService.GetTransferResult(ctx, info, token)
}

func (s *server) GetTripByCheapest(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	return s.travelPlanService.GetByCheapest(ctx, info, token)
}

func (s *server) GetTripByQuickest(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	return s.travelPlanService.GetByQuickest(ctx, info, token)
}

func (s *server) GetTripByMinStation(ctx context.Context, info util.RoutePlanInfo, token string) ([]util.TripDetails, error) {
	return s.travelPlanService.GetByMinStation(ctx, info, token)
}

func (s *server) GetRegionList(ctx context.Context, token string) (string, error) {
	return s.ticketOfficeService.GetRegionList(ctx, token)
}

func (s *server) GetAllOffices(ctx context.Context, token string) ([]util.Office, error) {
	return s.ticketOfficeService.GetAll(ctx, token)
}

func (s *server) GetSpecificOffices(ctx context.Context, province string, city string, region string, token string) ([]util.Office, error) {
	return s.ticketOfficeService.GetSpecificOffices(ctx, province, city, region, token)
}

func (s *server) AddOffice(ctx context.Context, province string, city string, region string, office util.Office, token string) (string, error) {
	return s.ticketOfficeService.AddOffice(ctx, province, city, region, token, office)
}

func (s *server) DeleteOffice(ctx context.Context, province string, city string, region string, officeName string, token string) (string, error) {
	return s.ticketOfficeService.DeleteOffice(ctx, province, city, region, officeName, token)
}

func (s *server) UpdateOffice(ctx context.Context, province string, city string, region string, oldOfficeName string, newOffice util.Office, token string) (string, error) {
	return s.ticketOfficeService.UpdateOffice(ctx, province, city, region, oldOfficeName, token, newOffice)
}

func (s *server) GetAllInsurances(ctx context.Context, token string) ([]util.Insurance, error) {
	return s.insuranceService.GetAllInsurances(ctx, token)
}

func (s *server) GetAllInsuranceTypes(ctx context.Context, token string) ([]util.InsuranceType, error) {
	return s.insuranceService.GetAllInsuranceTypes(ctx, token)
}

func (s *server) DeleteInsurance(ctx context.Context, insuranceId string, token string) (string, error) {
	return s.insuranceService.DeleteInsurance(ctx, insuranceId, token)
}

func (s *server) DeleteInsuranceByOrderId(ctx context.Context, orderID string, token string) (string, error) {
	return s.insuranceService.DeleteInsuranceByOrderId(ctx, orderID, token)
}

func (s *server) ModifyInsurance(ctx context.Context, insuranceId string, orderId string, typeIndex int, token string) (util.Insurance, error) {
	return s.insuranceService.ModifyInsurance(ctx, uint16(typeIndex), insuranceId, orderId, token)
}

func (s *server) CreateNewInsurance(ctx context.Context, typeIndex int, orderId string, token string) (util.Insurance, error) {
	return s.insuranceService.CreateNewInsurance(ctx, uint16(typeIndex), orderId, token)
}

func (s *server) GetInsuranceById(ctx context.Context, id string, token string) (util.Insurance, error) {
	return s.insuranceService.GetInsuranceById(ctx, id, token)
}

func (s *server) FindInsuranceByOrderId(ctx context.Context, orderId string, token string) (util.Insurance, error) {
	return s.insuranceService.FindInsuranceByOrderId(ctx, orderId, token)
}

func (s *server) CreateAndModifyRoute(ctx context.Context, id string, startStation string, endStation string, stationList []string, distanceList []uint16, token string) (util.Route, error) {
	route := util.Route{Id: id, StartStationId: startStation, TerminalStationId: endStation, Stations: stationList, Distances: distanceList}
	return s.routeService.CreateAndModifyRoute(ctx, route, token)
}

func (s *server) DeleteRoute(ctx context.Context, routeId string, token string) (string, error) {
	return s.routeService.DeleteRoute(ctx, routeId, token)
}

func (s *server) QueryRoutesById(ctx context.Context, routeId string, token string) (util.Route, error) {
	return s.routeService.QueryById(ctx, routeId, token)
}

func (s *server) QueryAllRoutes(ctx context.Context, token string) ([]util.Route, error) {
	return s.routeService.QueryAll(ctx, token)
}

func (s *server) QueryRoutesByStartAndTerminal(ctx context.Context, startId string, terminalId string, token string) ([]util.Route, error) {
	return s.routeService.QueryByStartAndTerminal(ctx, startId, terminalId, token)
}

func (s *server) PayBookingDifference(ctx context.Context, info util.RebookInfo, token string) (util.Order, error) {
	return s.rebookService.PayDifference(ctx, info, token)
}

func validateRebookParams(logger *slog.Logger, r *http.Request) (model.RebookInfo, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.RebookInfo{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	loginId := r.Form.Get("loginId")
	orderId := r.Form.Get("orderId")
	oldTripId := r.Form.Get("oldTripId")
	tripId := r.Form.Get("tripId")
	seatType := r.Form.Get("seatType")
	date := r.Form.Get("date")

	//TO-DO
	// validate mandatory fields

	seatTypeUint16, err := util.StringToUint16(seatType)
	if err != nil {
		logger.Error(err.Error())
		return model.RebookInfo{}, "", err
	}

	return model.RebookInfo{LoginId: loginId,
		OrderId:   orderId,
		OldTripId: oldTripId,
		TripId:    tripId,
		SeatType:  seatTypeFloat32,
		Date:      date}, token, nil
}

func (s *server) rebook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/rebook")

	info, token, err := validateRebookParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	order, err := s.rebookService.Get().Rebook(ctx, info, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("rebook successfully executed!", "order", order)
}

func (s *server) QueryStationsById(ctx context.Context, stationId string, token string) (string, error) {
	return s.stationService.QueryById(ctx, stationId, token)
}

func (s *server) QueryStations(ctx context.Context, token string) ([]util.Station, error) {
	return s.stationService.Query(ctx, token)
}

func (s *server) CreateStation(ctx context.Context, station util.Station, token string) (util.Station, error) {
	return s.stationService.Create(ctx, station, token)
}

func (s *server) UpdateStation(ctx context.Context, station util.Station, token string) (util.Station, error) {
	return s.stationService.Update(ctx, station, token)
}

func (s *server) DeleteStation(ctx context.Context, station util.Station, token string) (string, error) {
	return s.stationService.Delete(ctx, station, token)
}

func (s *server) QueryForStationId(ctx context.Context, stationName string, token string) (string, error) {
	return s.stationService.QueryForStationId(ctx, stationName, token)
}

func (s *server) QueryForStationIdBatch(ctx context.Context, stationNameList []string, token string) ([]string, error) {
	return s.stationService.QueryForIdBatch(ctx, stationNameList, token)
}

func (s *server) QueryForStationNameBatch(ctx context.Context, stationIdList []string, token string) ([]string, error) {
	return s.stationService.QueryForNameBatch(ctx, stationIdList, token)
}
*/

/*func validateCalculateRefundParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	orderId := r.Form.Get("orderId")

	//TO-DO
	// validate mandatory fields

	return token, orderId, nil
}

func (s *server) calculateRefund(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/calculateRefund")

	token, orderId, err := validateCalculateRefundParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	refund, err := s.cancelService.Get().Calculate(ctx, orderId, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("refund successfully calculated!", "refund", refund)
}*/

func validateCancelTicketParams(logger *slog.Logger, r *http.Request) (string, string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", "", err
	}
	// get params
	token := r.Form.Get("token")
	orderId := r.Form.Get("orderId")
	loginId := r.Form.Get("loginId")

	return orderId, loginId, token, nil
}

func (s *server) cancelTicket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/cancelTicket")

	orderId, loginId, token, err := validateCancelTicketParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	result, err := s.cancelService.Get().CancelTicket(ctx, orderId, loginId, token)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug(result)
	tt_metrics.TicketsCanceled.Inc()
}

func (s *server) consistencyWindow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/consistencyWindow")

	result, _ := s.cancelService.Get().GetConsistencyWindowValues(ctx)

	jsonData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type and write the JSON data
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *server) inconsistencies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/inconsistencies")

	result, _ := s.cancelService.Get().GetInconsistencies(ctx)

	w.Header().Set("Content-Type", "text/plain")
	resultStr := strconv.Itoa(result)

	w.Write([]byte(resultStr))
}

func (s *server) reset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/reset")

	s.cancelService.Get().Reset(ctx)

	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte("done!"))
}

/*func validatePayParams(logger *slog.Logger, r *http.Request) (string, string, string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", "", "", err
	}
	// get params
	token := r.Form.Get("token")
	orderId := r.Form.Get("orderId")
	price := r.Form.Get("price")
	userId := r.Form.Get("userId")

	//TO-DO
	// validate mandatory fields

	return orderId, price, userId, token, nil
}

func (s *server) pay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/pay")

	orderId, price, userId, token, err := validatePayParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	payment, err := s.paymentService.Get().Pay(ctx, orderId, price, userId, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("payment successfully executed!", "payment", payment)
}*/

/*func validateAddMoneyToAccountParams(logger *slog.Logger, r *http.Request) (string, string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", "", err
	}
	// get params
	token := r.Form.Get("token")
	price := r.Form.Get("price")
	userId := r.Form.Get("userId")

	//TO-DO
	// validate mandatory fields

	return userId, price, token, nil
}

func (s *server) addMoneyToAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/addMoneyToAccount")

	userId, price, token, err := validateAddMoneyToAccountParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	payment, err := s.paymentService.Get().AddMoney(ctx, userId, price, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("money successfully added!", "payment", payment)
}

func (s *server) QueryPayments(ctx context.Context, token string) ([]util.Payment, error) {
	return s.paymentService.Query(ctx, token)
}
func (s *server) PayInside(ctx context.Context, tripId string, userId string, orderId string, token string) (string, error) {
	return s.insidePaymentService.Pay(ctx, tripId, userId, orderId, token)
}
func (s *server) CreatePaymentAccount(ctx context.Context, money string, userId string, token string) (string, error) {
	return s.insidePaymentService.CreateAccount(ctx, money, userId, token)
}
func (s *server) QueryAccount(ctx context.Context, token string) ([]util.Balances, error) {
	return s.insidePaymentService.QueryAccount(ctx, token)
}
func (s *server) DrawBack(ctx context.Context, userId string, money string, token string) (string, error) {
	return s.insidePaymentService.DrawBack(ctx, userId, money, token)
}
func (s *server) PayDifference(ctx context.Context, orderId string, userId string, price string, token string) (string, error) {
	return s.insidePaymentService.PayDifference(ctx, orderId, userId, price, token)
}
func (s *server) QueryAddMoney(ctx context.Context, token string) ([]util.AddMoney, error) {
	return s.insidePaymentService.QueryAddMoney(ctx, token)
}
func (s *server) NotifyPreserveSuccess(ctx context.Context, info util.NotificationInfo, token string) error {
	return s.notificationService.PreserveSuccess(ctx, info, token)
}
func (s *server) NotifyOrderCreateSuccess(ctx context.Context, info util.NotificationInfo, token string) error {
	return s.notificationService.OrderCreateSuccess(ctx, info, token)
}
func (s *server) NotifyOrderChangedSuccess(ctx context.Context, info util.NotificationInfo, token string) error {
	return s.notificationService.OrderChangedSuccess(ctx, info, token)
}
func (s *server) NotifyOrderCanceledSuccess(ctx context.Context, info util.NotificationInfo, token string) error {
	return s.notificationService.OrderCancelSuccess(ctx, info, token)
}
func (s *server) BasicQueryForTravel(ctx context.Context, info util.Travel, token string) (util.TravelResult, error) {
	return s.basicService.QueryForTravel(ctx, info, token)
}
func (s *server) TicketInfoQueryForTravel(ctx context.Context, info util.Travel, token string) (util.TravelResult, error) {
	return s.ticketInfoService.QueryForTravel(ctx, info, token)
}
func (s *server) QueryPrice(ctx context.Context, routeId string, trainType string, token string) (util.PriceConfig, error) {
	return s.priceService.Query(ctx, routeId, trainType, token)
}
func (s *server) QueryAllPrices(ctx context.Context, token string) ([]util.PriceConfig, error) {
	return s.priceService.QueryAll(ctx, token)
}
func (s *server) CreatePriceConfig(ctx context.Context, info util.PriceConfig, token string) (util.PriceConfig, error) {
	return s.priceService.Create(ctx, info, token)
}
func (s *server) UpdatePriceConfig(ctx context.Context, info util.PriceConfig, token string) (util.PriceConfig, error) {
	return s.priceService.Update(ctx, info, token)
}
func (s *server) Preserve(ctx context.Context, info util.OrderInfo, token string) (string, error) {
	return s.preserveService.Preserve(ctx, info, token)
}
func (s *server) PreserveOther(ctx context.Context, info util.OrderInfo, token string) (string, error) {
	return s.preserveOtherService.Preserve(ctx, info, token)
}
func (s *server) RefreshOrders(ctx context.Context, orderInfo util.OrderInfo, accountId string, token string) ([]util.Order, error) {
	return s.orderService.QueryOrdersForRefresh(ctx, orderInfo, accountId, token)
}
func (s *server) RefreshOtherOrders(ctx context.Context, orderInfo util.OrderInfo, accountId string, token string) ([]util.Order, error) {
	return s.orderOtherService.QueryOrdersForRefresh(ctx, orderInfo, accountId, token)
}*/

/*func validateGetAllContactsParams(logger *slog.Logger, r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	// get params
	token := r.Form.Get("token")

	//TO-DO
	// validate mandatory fields

	return token, nil
}

func (s *server) getAllContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/getAllContacts")

	token, err := validateGetAllContactsParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.contactService.Get().GetAllContacts(ctx, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func validateCreateNewContactParams(logger *slog.Logger, r *http.Request) (model.Contact, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.Contact{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	accountID := r.Form.Get("accountID")
	name := r.Form.Get("name")
	documentType := r.Form.Get("documentType")
	documentNumber := r.Form.Get("documentNumber")
	phoneNumber := r.Form.Get("phoneNumber")

	//TO-DO
	// validate mandatory fields

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.Contact{}, "", err
	}

	return model.Contact{AccountId: accountID,
		Name:           name,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNumber,
		PhoneNumber:    phoneNumber}, token, nil
}

func (s *server) createNewContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/createNewContact")

	contact, token, err := validateCreateNewContactParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	new, err := s.contactService.Get().CreateNewContacts(ctx, contact, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("contact %s successfully added", new.Id)
}

func validateDeleteContactParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	contactId := r.Form.Get("contactId")

	//TO-DO
	// validate mandatory fields

	return contactId, token, nil
}

func (s *server) deleteContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/deleteContact")

	contactId, token, err := validateDeleteContactParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	result, err := s.contactService.Get().DeleteContacts(ctx, contactId, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug(result)
}

func validateModifyContactParams(logger *slog.Logger, r *http.Request) (model.Contact, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.Contact{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	contactId := r.Form.Get("contactId")
	accountID := r.Form.Get("accountID")
	name := r.Form.Get("name")
	documentType := r.Form.Get("documentType")
	documentNumber := r.Form.Get("documentNumber")
	phoneNumber := r.Form.Get("phoneNumber")

	//TO-DO
	// validate mandatory fields

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.Contact{}, "", err
	}

	return model.Contact{Id: contactId,
		AccountId:      accountID,
		Name:           name,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNumber,
		PhoneNumber:    phoneNumber}, token, nil
}

func (s *server) modifyContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/modifyContact")

	contact, token, err := validateModifyContactParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	contact, err = s.contactService.Get().ModifyContacts(ctx, contact, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func validateFindContactsByAccountIdParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	accountID := r.Form.Get("accountID")

	//TO-DO
	// validate mandatory fields

	return accountID, token, nil
}

func (s *server) findContactsByAccountId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/findContactsByAccountId")

	accountId, token, err := validateFindContactsByAccountIdParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.contactService.Get().FindContactsByAccountId(ctx, accountId, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

}

func validateGetContactsByContactIdParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	contactId := r.Form.Get("contactId")

	//TO-DO
	// validate mandatory fields

	return contactId, token, nil
}

func (s *server) getContactsByContactId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/getContactsByContactId")

	id, token, err := validateGetContactsByContactIdParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.contactService.Get().GetContactsByContactId(ctx, id, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}*/

/*
func (s *server) ExecuteTicket(ctx context.Context, orderId string, token string) (string, error) {
	return s.executeService.ExecuteTicket(ctx, orderId, token)
}

func (s *server) CollectTicket(ctx context.Context, orderId string, token string) (string, error) {
	return s.executeService.CollectTicket(ctx, orderId, token)
}

func (s *server) FindAllSecurityConfigs(ctx context.Context, token string) ([]util.SecurityConfig, error) {
	return s.securityService.FindAllSecurityConfigs(ctx, token)
}

func (s *server) CreateSecurityConfig(ctx context.Context, name string, value string, description string, token string) (util.SecurityConfig, error) {
	return s.securityService.Create(ctx, name, value, description, token)
}

func (s *server) UpdateSecurityConfig(ctx context.Context, id string, name string, value string, description string, token string) (util.SecurityConfig, error) {
	return s.securityService.Update(ctx, id, name, value, description, token)
}

func (s *server) DeleteSecurityConfig(ctx context.Context, id string, token string) (string, error) {
	return s.securityService.Delete(ctx, id, token)
}

func (s *server) PerformSecurityCheck(ctx context.Context, accountID string, token string) (string, error) {
	return s.securityService.Check(ctx, accountID, token)
}

func (s *server) QueryAllConfigs(ctx context.Context) ([]util.Config, error) {
	return s.configService.QueryAll(ctx)
}

func (s *server) CreateConfig(ctx context.Context, info util.Config) (util.Config, error) {
	return s.configService.CreateConfig(ctx, info)
}

func (s *server) UpdateConfig(ctx context.Context, info util.Config) (util.Config, error) {
	return s.configService.UpdateConfig(ctx, info)
}

func (s *server) DeleteConfig(ctx context.Context, configName string) (string, error) {
	return s.configService.DeleteConfig(ctx, configName)
}

func (s *server) RetrieveConfig(ctx context.Context, configName string) (util.Config, error) {
	return s.configService.Retrieve(ctx, configName)
}

func (s *server) CreateTrainType(ctx context.Context, id string, economyClass int, comfortClass int, avgSpeed int, token string) (util.Train, error) {
	train := util.Train{Id: id, EconomyClass: uint16(economyClass), ComfortClass: uint16(comfortClass), AvgSpeed: uint16(avgSpeed)}
	return s.trainService.Create(ctx, train, token)
}

func (s *server) UpdateTrainType(ctx context.Context, id string, economyClass int, comfortClass int, avgSpeed int, token string) (util.Train, error) {
	train := util.Train{Id: id, EconomyClass: uint16(economyClass), ComfortClass: uint16(comfortClass), AvgSpeed: uint16(avgSpeed)}
	return s.trainService.Update(ctx, train, token)
}

func (s *server) DeleteTrainType(ctx context.Context, id string, token string) (string, error) {
	return s.trainService.Delete(ctx, id, token)
}

func (s *server) QueryTrainTypes(ctx context.Context, token string) ([]util.Train, error) {
	return s.trainService.Query(ctx, token)
}

func (s *server) RetrieveTrainType(ctx context.Context, id string, token string) (util.Train, error) {
	return s.trainService.Retrieve(ctx, id, token)
}
*/

/*func validateGetAllUsersParams(logger *slog.Logger, r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", err
	}
	// get params
	token := r.Form.Get("token")

	//TO-DO
	// validate mandatory fields

	return token, nil
}

func (s *server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/getAllUsers")

	token, err := validateGetAllUsersParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.userService.Get().GetAllUsers(ctx, token)
	if err != nil {
		logger.Error(err.Error())
	}
}

func validateGetUserByIdParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	userID := r.Form.Get("userID")

	//TO-DO
	// validate mandatory fields

	return userID, token, nil
}

func (s *server) getUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/getUserById")

	userId, token, err := validateGetUserByIdParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.userService.Get().GetUserById(ctx, userId, token)
	if err != nil {
		logger.Error(err.Error())
	}
}

func validateGetUserByUsernameParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	username := r.Form.Get("username")

	//TO-DO
	// validate mandatory fields

	return username, token, nil
}

func (s *server) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/getUserByUsername")

	username, token, err := validateGetUserByUsernameParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.userService.Get().GetUserByUsername(ctx, username, token)
	if err != nil {
		logger.Error(err.Error())
	}
}

func validateDeleteUserByIdParams(logger *slog.Logger, r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}
	// get params
	token := r.Form.Get("token")
	userId := r.Form.Get("userId")

	//TO-DO
	// validate mandatory fields

	return userId, token, nil
}

func (s *server) deleteUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/deleteUserById")

	userId, token, err := validateDeleteUserByIdParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = s.userService.Get().DeleteUserById(ctx, userId, token)
	if err != nil {
		logger.Error(err.Error())
	}
}

func validateUpdateUserParams(logger *slog.Logger, r *http.Request) (model.User, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.User{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	gender := r.Form.Get("gender")
	documentType := r.Form.Get("documentType")
	documentNum := r.Form.Get("documentNum")
	email := r.Form.Get("email")

	//TO-DO
	// validate mandatory fields

	genderUint16, err := util.StringToUint16(gender)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	return model.User{Username: username,
		Password:       password,
		Gender:         genderUint16,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNum,
		Email:          email}, token, nil
}

func (s *server) updateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/updateUser")

	user, token, err := validateUpdateUserParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	user, err = s.userService.Get().UpdateUser(ctx, user, token)
	if err != nil {
		logger.Error(err.Error())
	}
}

func validateRegisterUserParams(logger *slog.Logger, r *http.Request) (model.User, string, error) {
	if err := r.ParseForm(); err != nil {
		return model.User{}, "", err
	}
	// get params
	token := r.Form.Get("token")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	gender := r.Form.Get("gender")
	documentType := r.Form.Get("documentType")
	documentNum := r.Form.Get("documentNum")
	email := r.Form.Get("email")

	//TO-DO
	// validate mandatory fields

	genderUint16, err := util.StringToUint16(gender)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	documentTypeUint16, err := util.StringToUint16(documentType)
	if err != nil {
		logger.Error(err.Error())
		return model.User{}, "", err
	}

	return model.User{Username: username,
		Password:       password,
		Gender:         genderUint16,
		DocumentType:   documentTypeUint16,
		DocumentNumber: documentNum,
		Email:          email}, token, nil
}

func (s *server) registerUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/registerUser")

	user, token, err := validateRegisterUserParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	user, err = s.userService.Get().RegisterUser(ctx, user, token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Debug("success! new registered user!", "username", user.Username)
	w.Header().Set("Content-Type", "application/json")
	userJson, err := json.Marshal(user)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	w.Write(userJson)
}*/

/*func (s *server) GetAllUsersAuth(ctx context.Context, token string) ([]util.User, error) {
	return s.authService.GetAllUsers(ctx, token)
}
func (s *server) DeleteUserByIdAuth(ctx context.Context, userId string, token string) (string, error) {
	return s.authService.DeleteUserById(ctx, userId, token)
}
func (s *server) CreateDefaultUser(ctx context.Context, username string, userId string, password string, token string) (string, error) {
	user := util.User{Username: username, Password: password, UserId: userId}
	return s.authService.CreateDefaultUser(ctx, user, token)
}*/

/*func validateLoginParams(logger *slog.Logger, r *http.Request) (string, string, string, model.Captcha, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", "", model.Captcha{}, err
	}
	// get params
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	verificationCode := r.Form.Get("verificationCode")

	//TO-DO
	//Get Captcha info

	//TO-DO
	// validate mandatory fields

	return username, password, verificationCode, model.Captcha{}, nil
}

type Token struct {
	Token string
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering /wrk2-api/user/login")

	username, password, verificationCode, captcha, err := validateLoginParams(logger, r)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	token, err := s.authService.Get().Login(ctx, username, password, verificationCode, captcha)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Debug("login executed successfully", "username", username, "token", token)
	w.Header().Set("Content-Type", "application/json")
	tokenJson, err := json.Marshal(Token{Token: token})
	if err != nil {
		logger.Error(err.Error())
		return
	}
	w.Write(tokenJson)
}*/

/*func (s *server) GenerateCaptcha(ctx context.Context, captcha util.Captcha) (util.Captcha, string, error) {
	return s.verificationCodeService.Generate(ctx, captcha)
}
func (s *server) GetAvailableTrips(ctx context.Context, startingPlace string, endPlace string, departureTime string, token string) ([]util.TripDetails, error) {
	return s.travelService.QueryInfo(ctx, startingPlace, endPlace, departureTime, token)
}
func (s *server) GetAvailableTripsOther(ctx context.Context, startingPlace string, endPlace string, departureTime string, token string) ([]util.TripDetails, error) {
	return s.travel2Service.QueryInfo(ctx, startingPlace, endPlace, departureTime, token)
}*/
