FROM gcr.io/distroless/static-debian12:nonroot

ARG TARGETPLATFORM
USER 20000:20000
COPY --chmod=555 ${TARGETPLATFORM}/external-dns-anexia-webhook /opt/external-dns-anexia-webhook/app

ENTRYPOINT ["/opt/external-dns-anexia-webhook/app"]
