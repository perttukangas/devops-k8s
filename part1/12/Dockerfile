FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

RUN apk add --no-cache ca-certificates \
    && update-ca-certificates

FROM scratch AS runner

WORKDIR /usr/src/app

COPY --from=builder /app/main .
COPY --from=builder /app/ui /usr/src/app/ui
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV GIN_MODE=release

CMD ["./main"]