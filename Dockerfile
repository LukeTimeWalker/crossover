FROM gcr.io/distroless/static-debian13

ARG TARGETOS
ARG TARGETARCH

COPY ${TARGETOS}/${TARGETARCH}/crossover /crossover

ENTRYPOINT ["/crossover"]
