package omq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	pmodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type MetricsSlice []*MetricsEnvelope

type MetricsEnvelope struct {
	Name    string    `json:"name"`
	Help    string    `json:"help,omitempty"`
	Type    string    `json:"type"`
	Metrics []*Metric `json:"metrics,omitempty"`
}

type Metric struct {
	Type   string            `json:"type"` // for JSON unmarshalling
	Labels map[string]string `json:"labels,omitempty"`
	Data   ValueMetricData   `json:"data"`
}

type metricForUnmarshal struct {
	Type   string            `json:"type"` // for JSON unmarshalling
	Labels map[string]string `json:"labels,omitempty"`
	Data   json.RawMessage   `json:"data"`
}

type ValueMetricData interface {
	_dummy()
}

type ValueMetricDataProcessor interface {
	Process(m *pmodel.Metric) ValueMetricData
}

func (m MetricsSlice) MarshallOpenMetricsText() ([]byte, error) {
	buff := &bytes.Buffer{}

	enc := expfmt.NewEncoder(buff, expfmt.FmtText)
	for _, envelope := range m {
		typ, ok := pmodel.MetricType_value[envelope.Type]
		if !ok {
			return nil, fmt.Errorf(`invalid OpenMetrics type has given "%s"`, envelope.Type)
		}
		metricType := pmodel.MetricType(typ)

		var err error
		metrics := make([]*pmodel.Metric, len(envelope.Metrics))
		for i, metric := range envelope.Metrics {
			metrics[i], err = metric.toPromMetrics()
			if err != nil {
				return nil, fmt.Errorf("failed to convert the metrics: %w", err)
			}
		}

		err = enc.Encode(&pmodel.MetricFamily{
			Name:   &envelope.Name,
			Help:   &envelope.Help,
			Type:   &metricType,
			Metric: metrics,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to encode the given metrics to OpenMetrics text: %w", err)
		}
	}

	omt, err := io.ReadAll(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read the encoded OpenMetrics text from a buffer: %w", err)
	}
	return omt, nil
}

func (met *Metric) toLabelPairs() []*pmodel.LabelPair {
	labelPairs := make([]*pmodel.LabelPair, len(met.Labels))
	labelCount := 0
	for key, value := range met.Labels {
		// capture variables {{{
		key := key
		value := value
		/// }}}

		labelPairs[labelCount] = &pmodel.LabelPair{
			Name:  &key,
			Value: &value,
		}
		labelCount++
	}
	return labelPairs
}

func (met *Metric) toPromMetrics() (*pmodel.Metric, error) {
	m := &pmodel.Metric{
		Label: met.toLabelPairs(),
	}

	switch pmodel.MetricType(pmodel.MetricType_value[met.Type]) {
	case pmodel.MetricType_COUNTER:
		m.Counter = met.Data.(*CounterMetric).ToPromMetrics()
	case pmodel.MetricType_GAUGE:
		m.Gauge = met.Data.(*GaugeMetric).ToPromMetrics()
	case pmodel.MetricType_SUMMARY:
		var err error
		m.Summary, err = met.Data.(*SummaryMetric).ToPromMetrics()
		if err != nil {
			return nil, err
		}
	case pmodel.MetricType_UNTYPED:
		m.Untyped = met.Data.(*UntypedMetric).ToPromMetrics()
	case pmodel.MetricType_HISTOGRAM:
		var err error
		m.Histogram, err = met.Data.(*HistogramMetric).ToPromMetrics()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(`invalid metric type has given "%s"`, met.Type)
	}

	return m, nil
}

func (met *Metric) UnmarshalJSON(inputBytes []byte) error {
	var input metricForUnmarshal
	err := json.Unmarshal(inputBytes, &input)
	if err != nil {
		return err
	}

	var data ValueMetricData
	switch pmodel.MetricType(pmodel.MetricType_value[input.Type]) {
	case pmodel.MetricType_COUNTER:
		data = &CounterMetric{}
	case pmodel.MetricType_GAUGE:
		data = &GaugeMetric{}
	case pmodel.MetricType_SUMMARY:
		data = &SummaryMetric{}
	case pmodel.MetricType_UNTYPED:
		data = &UntypedMetric{}
	case pmodel.MetricType_HISTOGRAM:
		data = &HistogramMetric{}
	default:
		return fmt.Errorf(`unsupported OpenMetrics type has given "%s"`, input.Type)
	}

	err = json.Unmarshal(input.Data, &data)
	if err != nil {
		return err
	}

	met.Type = input.Type
	met.Labels = input.Labels
	met.Data = data

	return nil
}
