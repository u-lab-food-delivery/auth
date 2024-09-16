# Dockerfile
FROM golang:1.22.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/app .
EXPOSE 50051
CMD ["./app"]