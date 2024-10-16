FROM node:22-alpine3.20 AS clientbuilder

WORKDIR /app

COPY ./client .

RUN npm ci && \
    npm run build

FROM golang:1.23.2-alpine3.20 AS serverbuilder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

RUN apk add --no-cache ca-certificates \
    && update-ca-certificates

FROM scratch AS runner

WORKDIR /usr/src/app

COPY --from=serverBuilder /app/main .
COPY --from=serverBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=clientBuilder /app/dist ./dist

ENV GIN_MODE=release

CMD ["./main"]