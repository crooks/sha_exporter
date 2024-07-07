package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix string = "sha"
)

type prometheusMetrics struct {
	groupSHA   *prometheus.GaugeVec
	groupUsers *prometheus.GaugeVec
	fileSHA    *prometheus.GaugeVec
}

func addPrefix(s string) string {
	return fmt.Sprintf("%s_%s", prefix, s)
}

func initCollectors() *prometheusMetrics {
	defaultGroupLabels := []string{"group", "gid"}
	sha := new(prometheusMetrics)

	sha.groupSHA = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("group_users_conforms"),
			Help: "Does the SHA256 hash of the users field match the configured hash.",
		},
		defaultGroupLabels,
	)
	prometheus.MustRegister(sha.groupSHA)

	sha.groupUsers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("group_users_count"),
			Help: "Number of users in the specified group.",
		},
		defaultGroupLabels,
	)
	prometheus.MustRegister(sha.groupUsers)

	sha.fileSHA = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("file_conforms"),
			Help: "Does the SHA256 of a file match the configured hash.",
		},
		[]string{"name"},
	)
	prometheus.MustRegister(sha.fileSHA)

	return sha
}
