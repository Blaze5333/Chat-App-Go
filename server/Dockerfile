FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY .env .env
EXPOSE 8080
CMD ["./server"]