FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache make
RUN go install github.com/air-verse/air@latest

ENV PATH="/root/go/bin:${PATH}"

EXPOSE 8080

CMD ["air"]