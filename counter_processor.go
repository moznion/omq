package omq

import pmodel "github.com/prometheus/client_model/go"

type CounterMetric struct {
	Value float64 `json:"value,omitempty"`
}

func (m *CounterMetric) _dummy() {}

func (m *CounterMetric) ToPromMetrics() *pmodel.Counter {
	return &pmodel.Counter{
		Value:    &m.Value,
		Exemplar: nil, // TODO
	}
}

type CounterProcessor struct{}

func (p *CounterProcessor) Process(m *pmodel.Metric) ValueMetricData {
	return &CounterMetric{
		Value: *m.Counter.Value,
	}
}
