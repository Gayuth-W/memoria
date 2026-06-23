FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o memoria main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/memoria .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./memoria"]
