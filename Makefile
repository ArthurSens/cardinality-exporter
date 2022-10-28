# Image URL to use all building/pushing image targets
IMG ?= ghcr.io/arthursens/cardinality-exporter:main

build: cardinality-exporter

# Build api binary
cardinality-exporter: fmt vet
	CGO_ENABLED=0 go build -v -ldflags '-w -extldflags '-static'' -o cardinality-exporter

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

docker-build:
	docker build . -t ${IMG}