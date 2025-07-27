FROM golang:1.24 as builder
WORKDIR /app
COPY go.* /app
RUN go mod download
COPY . .

RUN mkdir -p /app/bin
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -a \
  -ldflags "-s -w -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo 'v0.0.0')" \
  -o ./bin/ ./...

# Generate model.gob
RUN /app/bin/cli model:train

FROM scratch as runtime 
WORKDIR /app
COPY --from=builder /app/model.gob /app/bin/model.gob
COPY --from=builder /app/bin/* /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 443

CMD ["/app/telegram"]