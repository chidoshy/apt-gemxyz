FROM golang:1.18 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN GOOS=linux go build -o server

FROM ubuntu:20.04

# fix Kafka SASL connection
RUN apt-get update && apt-get install -y ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/server .
COPY ./migrations ./migrations
COPY ./data ./data
CMD ./server server
