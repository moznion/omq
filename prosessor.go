package omq

import (
	"encoding/json"
	"io"
	"log"

	"github.com/itchyny/gojq"
	pmodel "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func Query(reader io.ReadCloser, query string) ([]byte, error) {
	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return nil, err
	}

	jsonMetrics, err := ConvertOpenMetricsTextToJSON(reader)
	if err != nil {
		return nil, err
	}

	var metrics interface{}
	_ = json.Unmarshal(jsonMetrics, &metrics)

	var queried []byte
	iter := jqQuery.Run(metrics)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			log.Fatal(err)
		}
		marshal, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		queried = append(queried, marshal...)
	}

	return queried, nil
}

func toName2Metrics(reader io.ReadCloser) (MetricsSlice, error) {
	parser := &expfmt.TextParser{}
	parsed, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, err
	}

	ms := make(MetricsSlice, len(parsed))
	cnt := 0
	for k, v := range parsed {
		processor := getMetricsProcessor(v)
		if processor == nil {
			continue
		}

		metrics := make([]*Metric, len(v.Metric))
		for i, m := range v.Metric {
			labels := map[string]string{}
			for _, label := range m.Label {
				labels[*label.Name] = *label.Value
			}
			metrics[i] = &Metric{
				Type:   v.Type.String(),
				Labels: labels,
				Data:   processor.Process(m),
			}
		}

		ms[cnt] = &MetricsEnvelope{
			Name: k,
			Help: func() string {
				if v.Help == nil {
					return ""
				}
				return *v.Help
			}(),
			Type:    v.Type.String(),
			Metrics: metrics,
		}
		cnt++
	}

	return ms, nil
}

func ConvertOpenMetricsTextToJSON(reader io.ReadCloser) ([]byte, error) {
	mm, err := toName2Metrics(reader)
	if err != nil {
		return nil, err
	}

	marshal, err := json.Marshal(mm)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func ConvertJSONToOpenMetricsText(reader io.ReadCloser) ([]byte, error) {
	inputBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return ConvertJSONBytesToOpenMetricsText(inputBytes)
}

func ConvertJSONBytesToOpenMetricsText(inputBytes []byte) ([]byte, error) {
	var metricsSlice MetricsSlice
	err := json.Unmarshal(inputBytes, &metricsSlice)
	if err != nil {
		return nil, err
	}
	omt, err := metricsSlice.MarshallOpenMetricsText()
	if err != nil {
		return nil, err
	}
	return omt, nil
}

func getMetricsProcessor(v *pmodel.MetricFamily) ValueMetricDataProcessor {
	switch *v.Type {
	case pmodel.MetricType_COUNTER:
		return &CounterProcessor{}
	case pmodel.MetricType_GAUGE:
		return &GaugeProcessor{}
	case pmodel.MetricType_SUMMARY:
		return &SummaryProcessor{}
	case pmodel.MetricType_UNTYPED:
		return &UntypedProcessor{}
	case pmodel.MetricType_HISTOGRAM:
		return &HistogramProcessor{}
	default:
		return nil
	}
}
