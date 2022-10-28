package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace = "cardinality_exporter"
)

type Metrics struct {
	SeriesCountByMetricName     prometheus.GaugeVec
	LabelValueCountByLabelName  prometheus.GaugeVec
	MemoryInBytesByLabelName    prometheus.GaugeVec
	SeriesCountByLabelValuePair prometheus.GaugeVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		SeriesCountByMetricName: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "series_count_by_metric_name",
				Help:      "Timeseries count by metric name.",
			},
			[]string{"metric"},
		),
		LabelValueCountByLabelName: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "label_value_count_by_label_name",
				Help:      "Label values count by label.",
			},
			[]string{"label"},
		),
		MemoryInBytesByLabelName: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "memory_by_label_bytes",
				Help:      "Amount of memory used per label.",
			},
			[]string{"label"},
		),
		SeriesCountByLabelValuePair: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "series_count_by_label_value_pair",
				Help:      "Count of unique label/value pairs.",
			},
			[]string{"label"},
		),
	}
}

type TSDBStatus struct {
	Status string   `json:"status"`
	Data   TSDBData `json:"data"`
}

type TSDBData struct {
	SeriesCountByMetricName     []labelValuePair `json:"seriesCountByMetricName"`
	LabelValueCountByLabelName  []labelValuePair `json:"labelValueCountByLabelName"`
	MemoryInBytesByLabelName    []labelValuePair `json:"memoryInBytesByLabelName"`
	SeriesCountByLabelValuePair []labelValuePair `json:"seriesCountByLabelValuePair"`
}

type labelValuePair struct {
	Label string `json:"name"`
	Value uint64 `json:"value"`
}

func (m *Metrics) ProbeTSDBAPI(probeURL string, timeout time.Duration, logger log.Logger) error {

	apiURL := fmt.Sprintf("%s/api/v1/status/tsdb", probeURL)
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("Cannot create GET request to %v: %v", apiURL, err)
	}

	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Can't connect to %v: %v ", apiURL, err)
	}
	defer res.Body.Close()
	// Check the response and either log it, if 2xx, or return an error
	responseStatusLog := fmt.Sprintf("Request to %s returned status %s.", apiURL, res.Status)
	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		return errors.New(responseStatusLog)
	}
	level.Info(logger).Log(responseStatusLog)

	// Read the body of the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Can't read from socket: %v", err)
	}

	var tsdbStatus TSDBStatus
	// Parse the JSON response body into a struct
	err = json.Unmarshal(body, &tsdbStatus)
	if err != nil {
		return fmt.Errorf("Can't parse json: %v", err)
	}

	for _, labelValuePair := range tsdbStatus.Data.SeriesCountByMetricName {
		m.SeriesCountByMetricName.WithLabelValues(labelValuePair.Label).Set(float64(labelValuePair.Value))
	}

	for _, labelValuePair := range tsdbStatus.Data.LabelValueCountByLabelName {
		m.LabelValueCountByLabelName.WithLabelValues(labelValuePair.Label).Set(float64(labelValuePair.Value))
	}

	for _, labelValuePair := range tsdbStatus.Data.MemoryInBytesByLabelName {
		m.MemoryInBytesByLabelName.WithLabelValues(labelValuePair.Label).Set(float64(labelValuePair.Value))
	}

	for _, labelValuePair := range tsdbStatus.Data.SeriesCountByLabelValuePair {
		m.SeriesCountByLabelValuePair.WithLabelValues(labelValuePair.Label).Set(float64(labelValuePair.Value))
	}

	return nil
}
