package exporter

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpServer struct {
	mux *http.ServeMux
}

func NewHttpServer() *HttpServer {
	s := &HttpServer{
		mux: http.NewServeMux(),
	}

	s.mux.HandleFunc("/scrape", s.ScrapeHandler)
	return s
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//https://github.com/oliver006/redis_exporter/blob/master/exporter.go
func (s *HttpServer) ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		//e.targetScrapeRequestErrors.Inc()
		return
	}

	if !strings.Contains(target, "://") {
		target = "miio://" + target
	}

	u, err := url.Parse(target)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid 'target' parameter, parse err: %ck ", err), 400)
		//e.targetScrapeRequestErrors.Inc()
		return
	}

	token := u.Query().Get("token")

	registry := prometheus.NewRegistry()
	e, err := NewExporter(&ExporterTarget{
		Host:  u.Hostname(),
		Token: token,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid exporter target: %ck ", err), 400)
		//e.targetScrapeRequestErrors.Inc()
		return
	}
	registry.MustRegister(e)

	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
