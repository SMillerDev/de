FROM curlimages/curl:latest AS curl

# Create a production stage to run the application binary
FROM scratch

WORKDIR /app

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/de /app/de

COPY --from=curl /usr/bin/curl /usr/bin/curl


ENV PORT="8085"
ENV DATADIR="/data"
ENV DB_SEED="{}"

ENV GEO_BASE_JSON=
ENV ADMIN_USER=
ENV ADMIN_PASSWORD=

EXPOSE ${PORT}
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD ["/usr/bin/curl", "--fail", "--silent", "http://localhost:${PORT}/health"]
CMD ["/app/de"]
