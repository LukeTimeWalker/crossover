FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /crossover


FROM alpine:latest

WORKDIR /

COPY --from=build /crossover /crossover

ENTRYPOINT ["/crossover"]
