FROM golang:1.22-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o allora-producer ./cmd/producer

COPY ./config/config.example.yaml ./config.yaml

CMD ["./allora-producer"]

