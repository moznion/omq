package omq

import (
	"fmt"
	"strconv"

	pmodel "github.com/prometheus/client_model/go"
)

type HistogramMetric struct {
	Values      map[string]uint64 `json:"values,omitempty"`
	SampleSum   *float64          `json:"sampleSum,omitempty"`
	SampleCount *uint64           `json:"sampleCount,omitempty"`
}

func (m *HistogramMetric) _dummy() {}

func (m *HistogramMetric) ToPromMetrics() (*pmodel.Histogram, error) {
	buckets := make([]*pmodel.Bucket, len(m.Values))
	cnt := 0
	for upperBoundStr, cumulativeCount := range m.Values {
		cumulativeCount := cumulativeCount // capture
		upperBound, err := strconv.ParseFloat(upperBoundStr, 64)
		if err != nil {
			return nil, fmt.Errorf(`invalid upper bound valid has given "%s": %w`, upperBoundStr, err)
		}

		buckets[cnt] = &pmodel.Bucket{
			UpperBound:      &upperBound,
			CumulativeCount: &cumulativeCount,
			Exemplar:        nil, // TODO
		}
		cnt++
	}

	return &pmodel.Histogram{
		SampleCount: m.SampleCount,
		SampleSum:   m.SampleSum,
		Bucket:      buckets,
	}, nil
}

type HistogramProcessor struct{}

func (p *HistogramProcessor) Process(m *pmodel.Metric) ValueMetricData {
	histogram := m.Histogram

	values := map[string]uint64{}
	for _, bucket := range histogram.Bucket {
		upperBound := bucket.UpperBound
		if upperBound == nil {
			continue
		}

		values[fmt.Sprintf("%f", *upperBound)] = func() uint64 {
			cnt := bucket.CumulativeCount
			if cnt == nil {
				return 0
			}
			return *cnt
		}()
	}

	return &HistogramMetric{
		Values:      values,
		SampleCount: histogram.SampleCount,
		SampleSum:   histogram.SampleSum,
	}
}
