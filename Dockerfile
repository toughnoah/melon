FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY webhook .
USER 65532:65532

ENTRYPOINT ["/webhook"]
