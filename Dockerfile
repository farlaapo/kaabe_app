# Stage 1: Build the Go app
FROM golang:1.24.2-alpine AS build

WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY internal/config/config.yaml config/

RUN go build -o /kaabe ./cmd/server

# Stage 2: Run the Go app
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=build /kaabe /kaabe

#  Add wait-for-it.sh
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /kaabe /wait-for-it.sh

EXPOSE 8080

#  Wait for both db and redis before starting the app
CMD ["/bin/sh", "-c", "echo 'Starting app'; /wait-for-it.sh db:5432 -- /wait-for-it.sh redis:6379 -- /kaabe"]
