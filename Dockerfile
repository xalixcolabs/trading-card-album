# syntax=docker/dockerfile:1

# Empaqueta el binario ya compilado en local (ver `make build`) en una imagen
# distroless. No compila nada: solo copia build/trading-card-album.
#
#   make build
#   docker build -t trading-card-album:latest .

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY build/trading-card-album /app/trading-card-album

ENV APP_PORT=8080
ENV DATABASE_URL=sqlite:data/trading-card-album.sqlite3

EXPOSE 8080

CMD ["/app/trading-card-album"]