FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/bidding_svc ./cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auction_svc ./cmd/auction_svc/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notification_svc ./cmd/notification_svc/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/gateway_svc ./cmd/gateway/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /bin/bidding_svc .
COPY --from=builder /bin/auction_svc .
COPY --from=builder /bin/notification_svc .
COPY --from=builder /bin/gateway_svc .

EXPOSE 8081 50051 50052

CMD ["/bin/sh"]