FROM golang:1.23 as builder
WORKDIR /app
COPY go.* /app
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -a \
  -ldflags "-s -w" \
  -o /app/bin/telegram /app/cmd/telegram/main.go

FROM scratch as runtime 
COPY --from=builder /app/bin/* /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/