FROM golang:latest

RUN mkdir /app

WORKDIR /app

ADD . /app

RUN go get github.com/gorilla/websocket