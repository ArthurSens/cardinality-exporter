package main

import (
	"net/http"
	"os"
	"time"

	cardinality "github.com/ArthurSens/cardinality-exporter/pkg"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
)

const (
	metricsPath    = "/metrics"
	listenAddress  = ":8080"
	probeURL       = "http://localhost:9090"
	interval       = 1
	timeoutSeconds = 60
)

var (
	logger log.Logger
)

func main() {
	metrics := cardinality.NewMetrics()
	prometheus.MustRegister(version.NewCollector(cardinality.Namespace))
	prometheus.MustRegister(metrics.SeriesCountByMetricName)
	prometheus.MustRegister(metrics.LabelValueCountByLabelName)
	prometheus.MustRegister(metrics.MemoryInBytesByLabelName)
	prometheus.MustRegister(metrics.SeriesCountByLabelValuePair)

	logger = promlog.New(&promlog.Config{})
	level.Info(logger).Log("msg", "Starting http-prober", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	// Run initial probe
	err := metrics.ProbeTSDBAPI(probeURL, timeoutSeconds*time.Second, logger)
	if err != nil {
		level.Info(logger).Log("msg", "Failed to probe Prometheus TSDB status", "err", err)
	}

	go func() {
		t := time.NewTicker(interval * time.Hour)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				err := metrics.ProbeTSDBAPI(probeURL, timeoutSeconds*time.Second, logger)
				if err != nil {
					level.Info(logger).Log("msg", "Failed to probe Prometheus TSDB status", "err", err)
				}
			}
		}
	}()

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Cardinality Exporter</title></head>
			<body>
			<h1>Cardinality Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	level.Info(logger).Log("msg", "Listening on", "address", listenAddress)
	server := &http.Server{Addr: listenAddress}

	if err := web.ListenAndServe(server, "", logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
