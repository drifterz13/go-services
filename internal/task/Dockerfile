## Build
FROM golang:1.16-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /task-grpc

## Deploy
FROM alpine:latest

WORKDIR /
COPY --from=build /task-grpc /task-grpc

EXPOSE 50051

ENTRYPOINT ["/task-grpc"]

