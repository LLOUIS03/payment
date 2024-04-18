# Start from the latest golang base image
FROM golang:alpine3.19 as builder

LABEL maintainer="Luis Louis <luis.louis.castro@gmail.com>"

ARG CGO_ENABLED=0
WORKDIR /payment

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build ./cmd/main.go

EXPOSE 8090
ENTRYPOINT [ "./build" ] 