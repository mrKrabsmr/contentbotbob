FROM golang:1.21.3

WORKDIR /pars

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

EXPOSE 8001


