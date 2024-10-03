FROM golang:1.22-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o allora-producer ./cmd/producer

FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    rm -rf /var/cache/apk/*

WORKDIR /home/appuser/

COPY --from=builder /app/allora-producer .
COPY --from=builder /app/config/config.yaml .

USER appuser

CMD ["./allora-producer"]
