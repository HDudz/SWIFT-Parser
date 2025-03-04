FROM golang:1.24-alpine AS builder

LABEL authors="HDudz"

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/data /root/data


EXPOSE 8080

CMD ["./main"]
