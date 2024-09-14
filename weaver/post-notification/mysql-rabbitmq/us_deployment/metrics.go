package main

import "github.com/ServiceWeaver/weaver/metrics"

var (
	inconsistencies = metrics.NewCounter(
		"sn_inconsistencies",
		"The number of times an cross-service inconsistency has occured",
	)
	notificationsReceived = metrics.NewCounter(
		"notificationsReceived",
		"The number of notifications received",
	)
	readPostDurationMs = metrics.NewHistogram(
		"sn_read_post_duration_ms",
		"Duration of read operation in milliseconds in the us region",
		metrics.NonNegativeBuckets,
	)
	queueDurationMs = metrics.NewHistogram(
		"sn_queue_duration_ms",
		"Duration of queue in milliseconds",
		metrics.NonNegativeBuckets,
	)
	consistencyWindow = metrics.NewHistogram(
		"sn_consistency_window_ms",
		"Time taken between the post write on master and the post read on the replica",
		metrics.NonNegativeBuckets,
	)
)
