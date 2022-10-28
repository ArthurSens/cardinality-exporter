FROM golang:1.19.2 as builder

WORKDIR /workspace

COPY ./*.go ./
COPY ./go.mod ./
COPY ./go.sum ./
COPY ./Makefile Makefile
RUN go mod download
RUN make build


FROM alpine:3.16.0
WORKDIR /
COPY --from=builder /workspace/cardinality-exporter /usr/bin/cardinality-exporter

ENTRYPOINT ["/usr/bin/cardinality-exporter"]