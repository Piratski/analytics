FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o analytics ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/analytics .

EXPOSE 8080

CMD ["./analytics"]
