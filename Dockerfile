FROM golang:1.22 as builder

WORKDIR /app

COPY ./go.mod /app/go.mod
RUN go mod download

COPY . /app

RUN go build -o bin/service cmd/api/main.go

FROM ubuntu:22.04

WORKDIR /

RUN apt update
RUN apt install -y ca-certificates
COPY --from=builder /app/bin/service /bin/service