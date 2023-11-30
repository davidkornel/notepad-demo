package monitoring

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricServer struct {
	logger  logr.Logger
	port    int
	metrics *Metrics
	reg     *prometheus.Registry
	srv     *http.Server
}

func NewMetricServer(l logr.Logger, port int) *MetricServer {
	metricServer := &MetricServer{
		logger:  l,
		port:    port,
		metrics: &Metrics{},
		reg:     prometheus.NewRegistry(),
		srv:     &http.Server{Addr: fmt.Sprintf(":%d", port)},
	}
	metricServer.metrics.Init(metricServer.reg)

	return metricServer
}

func (p *MetricServer) Start() {
	log := p.logger.WithName("MetricServer")
	log.Info("Starting MetricServer HTTP server")

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(p.reg, promhttp.HandlerOpts{Registry: p.reg}))

		err := p.srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("metrics server: normal shutdown")
		} else {
			log.Error(err, "error running prometheus HTTP server on", "addr",
				fmt.Sprintf(":%d", p.port))
		}
	}()
}

func (p *MetricServer) Shutdown(ctx context.Context) {
	if err := p.srv.Shutdown(ctx); err != nil {
		p.logger.WithName("MetricServerShutdown").Error(err, "error happened while shutting down metric server")
	}
}
func (p *MetricServer) GetMetrics() *Metrics {
	return p.metrics
}
