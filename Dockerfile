FROM golang:1.21.6-alpine3.19

WORKDIR /var/www
COPY . /var/www

RUN go mod download && go mod verify