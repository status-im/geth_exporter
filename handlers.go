package main

import (
	"io"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	promlog "github.com/prometheus/common/log"
)

const rootHTML = `<html>
	<head>
		<title>Status Node Exporter</title>
	</head>
	<body>
		<h1>Status Node Exporter</h1>
		<p><a href="` + metricsPath + `">Metrics</a></p>
		</body>
</html>`

func rootHandler(w http.ResponseWriter, r *http.Request) {
	writeBody(w, rootHTML)
}

func handleCollectError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadGateway)
	log.Println(err)
	writeBody(w, "Bad Gateway")
}

func writeBody(w io.Writer, body string) {
	if _, err := w.Write([]byte(body)); err != nil {
		log.Println(err)
	}
}

// func metricsHandler(c *collector) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		body, err := c.collect()
// 		if err != nil {
// 			handleCollectError(w, err)
// 			return
// 		}

// 		writeBody(w, body)
// 	}
// }

func metricsHandler2(w http.ResponseWriter, r *http.Request) {
	registry := newRegistry(*ipcPath)
	gatherers := prometheus.Gatherers{
		registry,
	}

	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(gatherers,
			promhttp.HandlerOpts{
				ErrorLog:      promlog.NewErrorLogger(),
				ErrorHandling: promhttp.ContinueOnError,
			}),
	)
	h.ServeHTTP(w, r)
}
