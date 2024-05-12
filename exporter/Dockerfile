FROM golang:1.21 AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o exporter

FROM gcr.io/distroless/static-debian12

WORKDIR /work
COPY --from=builder /build/exporter ./exporter

CMD [ "/work/exporter" ]
