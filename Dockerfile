FROM gcr.io/distroless/static-debian12:nonroot

USER 20000:20000
COPY --chmod=555 external-dns-anexia-webhook /opt/external-dns-anexia-webhook/app

ENTRYPOINT ["/opt/external-dns-anexia-webhook/app"]
