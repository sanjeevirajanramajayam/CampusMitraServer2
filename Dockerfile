# Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates && update-ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

# Final lightweight image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
