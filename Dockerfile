FROM golang:latest AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o main app/cmd/main.go

ENTRYPOINT  ["/app/main", "test1.txt"]
#ENTRYPOINT  ["/app/main", "test2.txt"]
#ENTRYPOINT  ["/app/main", "test3.txt"]
#ENTRYPOINT  ["/app/main", "test4.txt"]
#ENTRYPOINT  ["/app/main", "test5.txt"]
#ENTRYPOINT  ["/app/main", "test6.txt"]
#ENTRYPOINT  ["/app/main", "test7.txt"]
#ENTRYPOINT  ["/app/main", "test8.txt"]
#ENTRYPOINT  ["/app/main", "test9.txt"]
#ENTRYPOINT  ["/app/main", "test10.txt"]