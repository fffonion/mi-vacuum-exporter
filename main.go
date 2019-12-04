package main

import (
	"net/http"

	"github.com/prometheus/common/log"

	"github.com/fffonion/mi-vacuum-exporter/exporter"
)

func main() {
	s := exporter.NewHttpServer()
	log.Infoln("Accepting Prometheus Requests on :9234")
	log.Fatal(http.ListenAndServe(":9234", s))
}
