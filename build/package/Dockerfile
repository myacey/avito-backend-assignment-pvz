# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY configs/ ./configs/

RUN CGO_ENABLED=0 GOOS=linux go build -v -o app cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/configs /root/configs
COPY --from=builder /app/.env .

# export ports
# 8080 - HTTP server
# 3000 - gRPC server
# 9000 - metrics server (for prometheus)
EXPOSE 8080 3000 9000

CMD [ "./app" ]