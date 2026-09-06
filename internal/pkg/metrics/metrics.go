package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/lo"
)

type Metrics struct {
	Observatory  *Observatory
	ScrapeErrors prometheus.Counter
}

func NewMetrics(prefix string) *Metrics {
	return &Metrics{
		Observatory: NewObservatory(prefix),

		ScrapeErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_scrape_errors_total", prefix),
			Help: "Total number of failed Xray observatory scrapes."}),
	}
}

func (m *Metrics) MustRegister(registry *prometheus.Registry) {
	m.Observatory.MustRegister(registry)

	for _, collector := range m.ListCollectors() {
		registry.MustRegister(collector)
	}
}

func (m *Metrics) ListCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.ScrapeErrors,
	}
}

func (o *Observatory) SetAlive(labels LabelsObservatory, input bool) {
	o.Alive.With(labels.ToHash()).Set(float64(lo.Ternary(input, 1, 0)))
}

func (o *Observatory) SetDelay(labels LabelsObservatory, input int64) {
	o.Delay.With(labels.ToHash()).Set(float64(input))
}

func (o *Observatory) SetLastSeen(labels LabelsObservatory, input int64) {
	o.LastSeen.With(labels.ToHash()).Set(float64(input))
}

func (o *Observatory) SetLastTry(labels LabelsObservatory, input int64) {
	o.LastTry.With(labels.ToHash()).Set(float64(input))
}

func (m *Metrics) IncScrapeErrors() {
	m.ScrapeErrors.Inc()
}

func (m *Metrics) GetObservatory() IObservatory {
	return m.Observatory
}
