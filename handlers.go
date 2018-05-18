package main

import (
	"io"
	"log"
	"net/http"
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

func metricsHandler(c *collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := c.collect()
		if err != nil {
			handleCollectError(w, err)
			return
		}

		writeBody(w, body)
	}
}
