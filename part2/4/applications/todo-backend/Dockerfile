FROM golang:1.23.2-alpine3.20 AS serverbuilder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM scratch AS runner

WORKDIR /usr/src/app

COPY --from=serverBuilder /app/main .

ENV GIN_MODE=release

CMD ["./main"]