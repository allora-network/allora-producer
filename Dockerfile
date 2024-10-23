FROM golang:1.22-alpine
RUN apk add git

# Default to localnet
ARG CHAIN_NAME=allora-localnet-1

ENV GO111MODULE=on \
    CGO_ENABLED=0

RUN apk add --no-cache jq


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -buildvcs -o allora-producer ./cmd/producer

COPY ./config/config.${CHAIN_NAME}.yaml ./config.yaml

CMD ["./allora-producer"]
