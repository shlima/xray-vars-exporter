package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/lo"
)

var DefaultObservatoryDelayHistogramBuckets = []float64{
	1,
	2,
	5,
	10,
	20,
	30,
	40,
	50,
	75,
	100,
	150,
	200,
	300,
	500,
	1000,
	2000,
	5000,
	10000,
}

type Observatory struct {
	Alive     *prometheus.GaugeVec
	Delay     *prometheus.GaugeVec
	DelayHist *prometheus.HistogramVec
	LastSeen  *prometheus.GaugeVec
	LastTry   *prometheus.GaugeVec
}

func NewObservatory() *Observatory {
	lables := lo.Keys(NewLablesObservatory().ToHash())

	return &Observatory{
		Alive: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "xray_observatory_alive",
			Help: "Whether the Xray observatory outbound is alive.",
		},
			lables),

		Delay: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "xray_observatory_delay_ms",
				Help: "Xray observatory probe delay in milliseconds.",
			},
			lables),

		DelayHist: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "xray_observatory_delay_ms_histogram",
			Help:    "Xray observatory probe delay in milliseconds",
			Buckets: DefaultObservatoryDelayHistogramBuckets},
			lables),

		LastSeen: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "xray_observatory_last_seen_timestamp",
			Help: "Unix timestamp when the outbound was last seen alive."},
			lables),

		LastTry: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "xray_observatory_last_try_timestamp",
			Help: "Unix timestamp when the outbound was last probed."},
			lables),
	}
}

func (o *Observatory) ListCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		o.Alive,
		o.Delay,
		o.DelayHist,
		o.LastSeen,
		o.LastTry,
	}
}

func (o *Observatory) MustRegister(registry *prometheus.Registry) {
	for _, collector := range o.ListCollectors() {
		registry.MustRegister(collector)
	}
}
