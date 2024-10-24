### This is a multistage Dockerfile

## Build the application
FROM golang:1.23 AS build-stage

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownload them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy the source code
COPY . .

# set the env variable for CGO_ENABLED and CC
ENV CGO_ENABLED=1
ENV CC=gcc

# build the application
RUN go build -v -o taskmaster ./cmd/server/main.go

## Deploy the application into a lean image
FROM alpine:latest

WORKDIR /app

COPY --from=build-stage /app/taskmaster .

COPY --from=build-stage /app/.env .

# install sqlite3 database
RUN apk add --no-cache sqlite

EXPOSE 8000

# command to run the application
CMD ["./taskmaster"]
