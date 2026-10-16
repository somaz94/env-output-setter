# Build stage
# Pinned to the BUILD platform, not the target: the image is published
# multi-arch, and letting the toolchain stage run under QEMU emulation to
# produce an arm64 binary is minutes of emulated compilation for no reason. Go
# cross-compiles natively instead, driven by the TARGET* args buildx injects.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

WORKDIR /app

# BuildKit sets both for the target platform, the host's for a plain
# `docker build` without --platform; the TARGETOS default just restates linux.
#
# BuildKit is required either way — `$BUILDPLATFORM` on the FROM above is a
# BuildKit-only variable, and the legacy builder fails to parse the line rather
# than ignoring it. That is the default builder since Docker 23, so this only
# bites an explicit DOCKER_BUILDKIT=0.
ARG TARGETOS=linux
ARG TARGETARCH

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /env-output-setter ./cmd/main.go

# Final stage
FROM alpine:3.24

WORKDIR /app

COPY --from=builder /env-output-setter /env-output-setter

ENTRYPOINT ["/env-output-setter"]