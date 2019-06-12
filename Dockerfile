FROM golang:1.12

LABEL maintainer="Aidar Babanov <aidar.babanov@nu.edu.kz>"

ENV GO111MODULE=on mod=vendor CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOFLAGS=-mod=vendor

RUN mkdir -p app
WORKDIR src/app

COPY . .

RUN go build -o main

# EXPOSE 8080
ENTRYPOINT ["./main"]