FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls -la cmd > /ls_output.txt && cat /ls_output.txt
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/weather-service

EXPOSE 8080
CMD ["/app/weather-service"]
