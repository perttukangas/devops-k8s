FROM postgres:17.0-alpine3.20

WORKDIR /app

RUN apk add --no-cache curl

COPY . .

RUN chmod +x ./run.sh

CMD ["./run.sh"]