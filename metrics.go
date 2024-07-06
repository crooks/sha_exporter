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
}

func addPrefix(s string) string {
	return fmt.Sprintf("%s_%s", prefix, s)
}

func initCollectors() *prometheusMetrics {
	defaultLabels := []string{"group", "gid"}
	sha := new(prometheusMetrics)

	sha.groupSHA = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("group_users_conforms"),
			Help: "Does the SHA256 hash of the users field match the configured hash.",
		},
		defaultLabels,
	)
	prometheus.MustRegister(sha.groupSHA)

	sha.groupUsers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: addPrefix("group_users_count"),
			Help: "Number of users in the specified group.",
		},
		defaultLabels,
	)
	prometheus.MustRegister(sha.groupUsers)

	return sha
}
