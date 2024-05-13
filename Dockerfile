FROM golang:alpine AS builder

WORKDIR /usr/local/src

COPY app ./

RUN go build -o ./bin/app cmd/main.go

#CMD [". /main test1.txt"]