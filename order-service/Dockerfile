FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o order-service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/order-service .

EXPOSE 8081

CMD ["./order-service"]