package metrics

import "github.com/prometheus/client_golang/prometheus"

type IPrometheusDelete interface {
	Delete(labels prometheus.Labels) bool
}

type ICollector interface {
	ListCollectors() []prometheus.Collector
	MustRegister(registry *prometheus.Registry)
}

type IObservatory interface {
	ICollector
	SetAlive(labels LabelsObservatory, input bool)
	SetDelay(labels LabelsObservatory, input int64)
	SetLastSeen(labels LabelsObservatory, input int64)
	SetLastTry(labels LabelsObservatory, input int64)
	ObserveDelayHist(labels LabelsObservatory, input int64)
	DeleteLabelValues(labels LabelsObservatory)
}

type IMetrics interface {
	ICollector
	GetObservatory() IObservatory
	IncScrapeErrors()
}
