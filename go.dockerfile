FROM golang:1.22.1

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o bin/accounts-backend ./cmd/api

EXPOSE 8002

CMD ["./bin/accounts-backend"]
