package metrics

import "github.com/ServiceWeaver/weaver/metrics"

type RegionLabel struct {
	Region string
}

var (
	// post storage service
	ReadPostDurationMs = metrics.NewHistogramMap[RegionLabel](
		"sn_read_post_duration_ms",
		"Duration of a read operation in milliseconds in the current region",
		[]float64{},
	)
	// write home timeline service
	QueueDurationMs = metrics.NewHistogramMap[RegionLabel](
		"sn_queue_duration_ms",
		"Duration of queue in milliseconds in the current region",
		[]float64{},
	)
	ReceivedNotifications = metrics.NewCounterMap[RegionLabel](
		"sn_received_notifications",
		"The number of received notifications in the current region",
	)
	Inconsistencies = metrics.NewCounterMap[RegionLabel](
		"sn_inconsistencies",
		"The number of times an cross-service inconsistency has occured in the current region",
	)
	UpdateInconsistencies = metrics.NewCounterMap[RegionLabel](
		"sn_update_inconsistencies",
		"The number of times an cross-service inconsistency has occured in the current region after an update operation",
	)
	ConsistencyWindow = metrics.NewHistogramMap[RegionLabel](
		"sn_consistency_window_ms",
		"Time taken between the post write on master and the post read on the replica",
		[]float64{},
	)
)
