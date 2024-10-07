FROM golang:1.22-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0

RUN addgroup -S appgroup && adduser -S -D -h /home/appuser appuser -G appgroup

WORKDIR /home/appuser/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o allora-producer ./cmd/producer

COPY ./config/config.example.yaml ./config.yaml

USER appuser

CMD ["./allora-producer"]
