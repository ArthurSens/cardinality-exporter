# cardinality-exporter

Inspired by [thought-machine/prometheus-cardinality-exporter](https://github.com/thought-machine/prometheus-cardinality-exporter), exposes metrics from TSDB-status. The mentioned project is built with tooling made in-house, making it difficult to use and contribute. Thus, this project was created with tools that are broadly used in the market.

The cardinality-exporter will hit the endpoint http://localhost:9090/api/v1/status/tsdb and expose the response as Prometheus metrics. Since the Prometheus address is hardcoded, it can only work as a sidecar at the moment. Contributions are welcome!

---

## Development

Build and run the binary locally with:

```
make build
./cardinality-exporter
```

Or as a docker image:
```
make docker-build
docker run -p 8080:8080 ghcr.io/arthursens/cardinality-exporter:main
```

By clicking in the button below, you'll have ready-to-go Cloud Develepor Environment to develop cardinality-exporter:

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/arthursens/cardinality-exporter)