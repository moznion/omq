package omq

import pmodel "github.com/prometheus/client_model/go"

type UntypedMetric struct {
	Value float64 `json:"value,omitempty"`
}

func (m *UntypedMetric) _dummy() {}

func (m *UntypedMetric) ToPromMetrics() *pmodel.Untyped {
	return &pmodel.Untyped{
		Value: &m.Value,
	}
}

type UntypedProcessor struct{}

func (p *UntypedProcessor) Process(m *pmodel.Metric) ValueMetricData {
	return &UntypedMetric{
		Value: *m.Untyped.Value,
	}
}
