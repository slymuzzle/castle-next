# ==================================================================================== #
# BASE IMAGE
# ==================================================================================== #
FROM golang:1.22.5-alpine AS base

WORKDIR /workspace

# https://megamorf.gitlab.io/2019/09/08/alpine-go-builds-with-cgo-enabled/
ENV CGO_ENABLED=1

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    apk add --no-cache \
    build-base=0.5-r3 \
    curl=8.10.0-r0 \
    && go mod download \
    && go install github.com/go-task/task/v3/cmd/task@latest

# ==================================================================================== #
# DEVELOPMENT IMAGE
# ==================================================================================== #
FROM base AS development

ARG USER_UID=1000
ARG USER_GID=${USER_UID}

# Set SHELL flags for RUN commands to allow -e and pipefail
# Rationale: https://github.com/hadolint/hadolint/wiki/DL4006
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

# Add group and user for dev environment
RUN addgroup --system --gid ${USER_GID} devgroup \
    && adduser --system --uid ${USER_UID} devuser \
    && chown -R devuser /go \
    && curl -sSf https://atlasgo.sh | ash

USER devuser

CMD ["task", "cleanup", "watch"]

# ==================================================================================== #
# BUILD IMAGE
# ==================================================================================== #
FROM base AS build

# Add non-root user
RUN adduser --system --uid 1001 nonroot

COPY . .

RUN task build

# ==================================================================================== #
# PRODUCTION IMAGE
# ==================================================================================== #
FROM alpine:3.20 AS server

WORKDIR /production

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /workspace/.tmp/server server
COPY --from=build /workspace/.tmp/seed seed
COPY --from=build /workspace/configs configs

USER nonroot

CMD ["/production/server"]

# ==================================================================================== #
# MIGRATIONS IMAGE
# ==================================================================================== #
FROM arigaio/atlas:0.25.0-alpine AS migrations

COPY migrations /migrations
