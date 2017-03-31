package main

// Current Cost CC128 Real-Time power usage exported as Prometheus metrics.

// Reads serial port data as per http://www.currentcost.com/cc128/xml.htm

import (
	"flag"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	httpAddr = flag.String("http", ":6799", "listen on this address")
)

func main() {
	flag.Parse()

	go func() {
		for {
			readInput()
			// readInput only returns in case of failure.
			// On failure, back off before retrying.
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "currentcost_exporter\n")
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("listening on %v", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
