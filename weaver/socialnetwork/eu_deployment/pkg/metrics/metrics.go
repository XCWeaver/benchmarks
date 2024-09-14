package metrics

import "github.com/ServiceWeaver/weaver/metrics"

type RegionLabel struct {
	Region string
}

var (
	// wrk2 api
	UpdatePostDuration = metrics.NewHistogramMap[RegionLabel](
		"sn_update_post_duration_ms",
		"Duration of compose post endpoint in milliseconds in the current region",
		[]float64{},
	)
	// composed post service
	ComposedPosts = metrics.NewCounterMap[RegionLabel](
		"sn_composed_posts",
		"The number of composed posts in the current region",
	)
	UpdatedPosts = metrics.NewCounterMap[RegionLabel](
		"sn_updated_posts",
		"The number of composed posts in the current region",
	)
	SentNotifications = metrics.NewCounterMap[RegionLabel](
		"sn_sent_notifications",
		"The number of sent notifications in the current region",
	)
	// post storage service
	WritePostDurationMs = metrics.NewHistogramMap[RegionLabel](
		"sn_write_post_duration_ms",
		"Duration of a write operation in milliseconds in the current region",
		[]float64{},
	)
	ReadPostDurationMs = metrics.NewHistogramMap[RegionLabel](
		"sn_read_post_duration_ms",
		"Duration of a read operation in milliseconds in the current region",
		[]float64{},
	)
	UpdatePostDurationMs = metrics.NewHistogramMap[RegionLabel](
		"sn_update_post_operation_duration_ms",
		"Duration of a update operation in milliseconds in the current region",
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
)
