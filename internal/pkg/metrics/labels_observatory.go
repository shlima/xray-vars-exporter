package metrics

type LabelsObservatory struct {
	Outbound string
}

func NewLablesObservatory() *LabelsObservatory {
	return new(LabelsObservatory)
}

func (l *LabelsObservatory) ToHash() map[string]string {
	return map[string]string{"outbound": l.Outbound}
}

func (o *Observatory) ObserveDelayHist(labels LabelsObservatory, input int64) {
	o.DelayHist.With(labels.ToHash()).Observe(float64(input))
}

// DeleteLabelValues deletes not existing outbounds
func (o *Observatory) DeleteLabelValues(labels LabelsObservatory) {
	for _, collector := range o.ListCollectors() {
		if obj, ok := collector.(IPrometheusDelete); ok {
			obj.Delete(labels.ToHash())
		}
	}
}
