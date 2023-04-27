FROM golang:1.19 as builder
WORKDIR /app
ADD .. /app

RUN --mount=type=cache,target=/root/.cache/go-build go build -o bin/plant-shutter cmd/server/main.go

FROM ubuntu as plant-shutter
WORKDIR /
COPY --from=builder /app/bin/plant-shutter /usr/local/bin
USER root
