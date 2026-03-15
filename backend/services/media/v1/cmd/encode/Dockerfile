# syntax=docker/dockerfile:1.7
#
# Build (from repo root):
#   docker build -f backend/services/media/v1/cmd/encode/Dockerfile -t chirp-encode:latest .
#

ARG GO_VERSION=1.25.6

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS builder
WORKDIR /src/backend

# Cache deps
COPY backend/go.mod backend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build
COPY backend/ ./
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/encode ./services/media/v1/cmd/encode

FROM debian:bookworm-slim AS runtime
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates \
      ffmpeg \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /out/encode /usr/local/bin/encode
RUN chown -R 10001:10001 /app
USER 10001:10001

ENTRYPOINT ["/usr/local/bin/encode"]
