FROM gcr.io/distroless/static-debian13

ARG TARGETOS
ARG TARGETARCH
ARG PROJECT_VERSION

COPY target/crossover-${PROJECT_VERSION}-${TARGETOS}-${TARGETARCH} /crossover

ENTRYPOINT ["/crossover"]
