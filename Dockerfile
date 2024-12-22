FROM gcr.io/distroless/static-debian12

COPY crossover /

ENTRYPOINT ["/crossover"]
