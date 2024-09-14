package metrics

import "github.com/XCWeaver/xcweaver/metrics"

var (
	OrderTicketDuration = metrics.NewHistogram(
		"tt_order_ticket_duration_ms",
		"Duration of order endpoint in milliseconds",
		[]float64{},
	)
	Orders = metrics.NewCounter(
		"sn_orders",
		"The number of orders",
	)
	TicketsCanceled = metrics.NewCounter(
		"sn_tickets_canceled",
		"The number of tickets canceled",
	)
	Inconsistencies = metrics.NewCounter(
		"tt_inconsistencies",
		"The number of times an cross-service inconsistency has occured",
	)
	ConsistencyWindow = metrics.NewHistogram(
		"sn_consistency_window_ms",
		"Time taken between the post write on master and the post read on the replica",
		[]float64{},
	)
)
