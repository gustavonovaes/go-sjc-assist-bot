FROM golang:1.23 as builder
COPY go.* /app
RUN go mod download
COPY . .

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -a \
  # -ldflags "-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'v0.0.0')" \
  -o ./bin/ ./cmd/telegram

FROM scratch as runtime 
COPY --from=builder /app/bin/* /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/