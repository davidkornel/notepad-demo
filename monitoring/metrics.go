package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	promRequestLabels        []string
	promEndpointRequestTotal *prometheus.CounterVec
	promNotesCurrentActive   prometheus.Gauge
}

func (m *Metrics) Init(reg *prometheus.Registry) {
	m.promRequestLabels = []string{"req_type", "endpoint"}

	m.promEndpointRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Number of requests at an endpoint.",
	}, m.promRequestLabels)

	m.promNotesCurrentActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "notes_currently_active",
		Help: "Number of notes currently in the database",
	})
	reg.MustRegister(m.promEndpointRequestTotal)
	reg.MustRegister(m.promNotesCurrentActive)
}

func (m *Metrics) IncrementRequestTotal(endpoint string, method string) {
	m.promEndpointRequestTotal.WithLabelValues(method, endpoint).Inc()
}

func (m *Metrics) SetNumberOfActiveNotes(num int64) {
	m.promNotesCurrentActive.Set(float64(num))
}
