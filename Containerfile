FROM golang:alpine AS builder

WORKDIR /build

COPY . ./

RUN apk add --update gcc musl-dev sqlite-dev

RUN CGO_ENABLED=1 go build -o dist/de -tags "linux"

FROM alpine:latest

RUN apk add --update sqlite

WORKDIR /app

ARG TARGETPLATFORM
COPY --from=builder /build/dist/de /app/de

ENV PORT="8085"
ENV DATADIR="/data"
ENV DB_SEED="{}"

ENV GEO_BASE_JSON=
ENV ADMIN_USER=
ENV ADMIN_PASSWORD=

EXPOSE ${PORT}
HEALTHCHECK --interval=30s --timeout=3s --retries=3 CMD ["/usr/bin/curl", "--fail", "--silent", "http://localhost:${PORT}/health"]

RUN mkdir -p "${DATADIR}"
CMD ["/app/de"]
