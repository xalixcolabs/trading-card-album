# syntax=docker/dockerfile:1

# Empaqueta el binario ya compilado en local (ver `make build`) en una imagen
# distroless. No compila nada: solo copia build/trading-card-album.
# $ make build docker

FROM docker.io/debian:bookworm-slim

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY build/ .

ENV APP_PORT=8080
ENV DATABASE_URL=sqlite:database/trading-card-album.sqlite3

EXPOSE 8080

ENTRYPOINT ["/app/trading-card-album"]