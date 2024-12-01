package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Masterminds/log-go"
)

func servRoot(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(`<html>
		<head><title>SHA Exporter API Server</title></head>
		<body>
		<h1>SHA Exporter API Server</h1>
		<p><a href='/api/v1/metrics'>Metrics</a></p>
		</body>
		</html>`))
	if err != nil {
		log.Warnf("Error on returning home page: %s", err)
	}
}

func servMetrics(w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(cfg.Metrics)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s\n", string(data))
}

func apiServer() {
	http.HandleFunc("/", servRoot)
	http.HandleFunc("/api/v1/metrics", servMetrics)
	http.ListenAndServe(":8080", nil)
}
