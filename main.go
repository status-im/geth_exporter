package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	metricsPath = "/metrics"
)

var (
	ipcPathFlag = flag.String("ipc", "", "path to ipc file")
	hostFlag    = flag.String("host", "", "http server host")
	portFlag    = flag.Int("port", 9200, "http server port")
)

func usage() {
	flag.Usage()
	os.Exit(1)
}

func requiredFlag(f string) {
	log.Printf("flag -%s is required\n", f)
	usage()
}

func parseFlags() {
	flag.Parse()

	if *ipcPathFlag == "" {
		requiredFlag("ipc")
	}

	if flag.NArg() > 0 {
		log.Printf("Extra args in command line: %v", flag.Args())
		usage()
	}
}

func main() {
	parseFlags()

	http.HandleFunc(metricsPath, metricsHandler(*ipcPathFlag))
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/", rootHandler)

	listenAddress := fmt.Sprintf("%s:%d", *hostFlag, *portFlag)

	log.Println("Listening on", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}
