package main

import (
	"fmt"
	"gitlab/sha_exporter/config"
	"net/http"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/crooks/jlog"
	loglevel "github.com/crooks/log-go-level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cfg   *config.Config
	flags *config.Flags
	prom  *prometheusMetrics
)

func bool2Float(b bool) (boolFloat float64) {
	if b {
		boolFloat = 1
	} else {
		boolFloat = 0
	}
	return
}

func metricsCollector() {
	interval := time.Duration(cfg.ScrapeInterval) * time.Second
	log.Infof("Parsing group file %s at interval %d seconds", cfg.GroupFile, cfg.ScrapeInterval)
	for {
		// Process the groups file
		countGroupSuccess, countGroupFail, err := findGroups(cfg.GroupFile)
		if err != nil {
			log.Fatal(err)
		}
		countGroupTotal := countGroupSuccess + countGroupFail
		log.Debugf("Processed %d group entries: Success=%d, Fail=%d", countGroupTotal, countGroupSuccess, countGroupFail)
		// Process file hashes
		countSuccess, countFail, countMissing := iterFiles()
		countTotal := countSuccess + countFail + countMissing
		log.Debugf("Processed %d files: Success=%d, Fail=%d, NotFound=%d", countTotal, countSuccess, countFail, countMissing)
		time.Sleep(interval)
	}
}

func main() {
	var err error
	flags = config.ParseFlags()
	cfg, err = config.ParseConfig(flags.Config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	// Define logging level and method
	loglev, err := loglevel.ParseLevel(cfg.Logging.LevelStr)
	if err != nil {
		log.Fatalf("unable to set log level: %v", err)
	}
	if cfg.Logging.Journal && jlog.Enabled() {
		log.Current = jlog.NewJournal(loglev)
	} else {
		log.Current = log.StdLogger{Level: loglev}
	}

	prom = initCollectors()
	go metricsCollector()
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(`<html>
		<head><title>SHA Exporter</title></head>
		<body>
		<h1>SHA Exporter</h1>
		<p><a href='/metrics'>Metrics</a></p>
		</body>
		</html>`))
		if err != nil {
			log.Warnf("Error on returning home page: %s", err)
		}
	})
	exporter := fmt.Sprintf("%s:%d", cfg.Exporter.Address, cfg.Exporter.Port)
	err = http.ListenAndServe(exporter, nil)
	if err != nil {
		log.Fatalf("HTTP listener failed: %v", err)
	}
}
