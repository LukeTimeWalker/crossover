FROM gcr.io/distroless/static-debian12

COPY target/crossover /

ENTRYPOINT ["/crossover"]
