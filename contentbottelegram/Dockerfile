FROM golang:1.21.3

WORKDIR /tg-bot

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

EXPOSE 7000

