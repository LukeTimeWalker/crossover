FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN GOOS=linux go build -o /crossover


FROM gcr.io/distroless/static-debian12

WORKDIR /crossover

COPY --from=build /crossover /crossover/crossover

ENTRYPOINT ["./crossover"]
