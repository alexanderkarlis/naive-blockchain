FROM golang:alpine AS builder
RUN mkdir -p /usr/local/app/server
WORKDIR /usr/local/app/server
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o ./naive-blockchain .
