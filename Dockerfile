# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build server
RUN go build -o /server cmd/server/main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./server"]
