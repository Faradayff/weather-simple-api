FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o weather-simple-api cmd/server.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/weather-simple-api .
EXPOSE 8080
CMD ["./weather-simple-api"]
