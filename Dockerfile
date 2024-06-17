# Multi-stage build to generate custom k6 with extension
FROM golang:1.21.11-bookworm as builder

WORKDIR /go/src/go.k6.io/k6

ADD . .

RUN apt-get update -y && \
    apt-get install -y build-essential git

RUN go install go.k6.io/xk6/cmd/xk6@latest

RUN CGO_ENABLED=1 xk6 build \
    --with github.com/grafana/xk6-sql=. \
    --output /tmp/k6

FROM gcr.io/distroless/base-debian12

USER 12345:12345

COPY --from=builder /tmp/k6 /k6

ENTRYPOINT ["/k6"]
