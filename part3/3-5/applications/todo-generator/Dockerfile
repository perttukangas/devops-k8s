FROM alpine:3.20.3

WORKDIR /app

RUN apk add --no-cache curl jq

COPY . .

RUN chmod +x ./run.sh

CMD ["./run.sh"]