FROM golang:1.24-alpine

LABEL authors="HDudz"

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/api

EXPOSE 8080

CMD ["./main"]
