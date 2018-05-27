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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeBody(w, "OK")
}

func writeBody(w io.Writer, body string) {
	if _, err := w.Write([]byte(body)); err != nil {
		log.Println(err)
	}
}

func metricsHandler(ipcPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := r.URL.Query()["collect[]"]
		registry := newRegistry(ipcPath, filters)
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
}
