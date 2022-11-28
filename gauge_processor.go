package omq

import pmodel "github.com/prometheus/client_model/go"

type GaugeMetric struct {
	Value float64 `json:"value,omitempty"`
}

func (m *GaugeMetric) _dummy() {}

func (m *GaugeMetric) ToPromMetrics() *pmodel.Gauge {
	return &pmodel.Gauge{
		Value: &m.Value,
	}
}

type GaugeProcessor struct{}

func (p *GaugeProcessor) Process(m *pmodel.Metric) ValueMetricData {
	return &GaugeMetric{
		Value: *m.Gauge.Value,
	}
}
