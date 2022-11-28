package omq

import (
	"fmt"
	"strconv"

	pmodel "github.com/prometheus/client_model/go"
)

type SummaryMetric struct {
	Values      map[string]float64 `json:"values,omitempty"`
	SampleSum   *float64           `json:"sampleSum,omitempty"`
	SampleCount *uint64            `json:"sampleCount,omitempty"`
}

func (m *SummaryMetric) _dummy() {}

func (m *SummaryMetric) ToPromMetrics() (*pmodel.Summary, error) {
	quantiles := make([]*pmodel.Quantile, len(m.Values))
	cnt := 0
	for q, v := range m.Values {
		v := v // capture
		q, err := strconv.ParseFloat(q, 64)
		if err != nil {
			return nil, fmt.Errorf(`invalid summary value has given "%s": %w`, q, err)
		}

		quantiles[cnt] = &pmodel.Quantile{
			Quantile: &q,
			Value:    &v,
		}
		cnt++
	}
	return &pmodel.Summary{
		SampleCount: m.SampleCount,
		SampleSum:   m.SampleSum,
		Quantile:    quantiles,
	}, nil
}

type SummaryProcessor struct{}

func (p *SummaryProcessor) Process(m *pmodel.Metric) ValueMetricData {
	summary := m.Summary

	values := map[string]float64{}
	for _, q := range summary.Quantile {
		quantile := q.Quantile
		if quantile == nil {
			continue
		}

		values[fmt.Sprintf("%f", *quantile)] = func() float64 {
			if q.Value == nil {
				return 0
			}
			return *q.Value
		}()
	}

	return &SummaryMetric{
		Values:      values,
		SampleSum:   summary.SampleSum,
		SampleCount: summary.SampleCount,
	}
}
