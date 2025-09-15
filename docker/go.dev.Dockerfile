FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache make
RUN go install github.com/air-verse/air@v1.62.0
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.3
RUN go install github.com/mitranim/gow@latest


ENV PATH="/root/go/bin:${PATH}"
EXPOSE 8080

CMD ["air"]