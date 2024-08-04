FROM golang:1.20.6-bullseye as builder
RUN apt-get update -y && apt-get install -y file && apt-get install -y upx
COPY . /workspace
WORKDIR /workspace
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -tags=jsoniter -o manager .
RUN upx ./manager
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-s -w" -o health-check cmd/health-check.go
RUN upx ./health-check

FROM scratch

COPY --from=builder /workspace/manager ./manager
COPY --from=builder /workspace/health-check ./healthcheck
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./manager"]