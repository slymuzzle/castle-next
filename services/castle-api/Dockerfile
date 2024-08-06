# ==================================================================================== #
# BASE IMAGE
# ==================================================================================== #
FROM golang:1.22.5-alpine as base

# https://megamorf.gitlab.io/2019/09/08/alpine-go-builds-with-cgo-enabled/
ENV CGO_ENABLED=1

# Set SHELL flags for RUN commands to allow -e and pipefail
# Rationale: https://github.com/hadolint/hadolint/wiki/DL4006
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

RUN apk add --no-cache \
    build-base=0.5-r3 \
    curl=8.9.0-r0 \
    make=4.4.1-r2 \
    # Install atlas
    && curl -sSf https://atlasgo.sh | ash

# ==================================================================================== #
# DEVELOPMENT IMAGE
# ==================================================================================== #
FROM base as development

WORKDIR /workspace

USER root

ARG USER_UID=1000
ARG USER_GID=${USER_UID}

# Add group and user for dev environment
RUN addgroup --system --gid ${USER_GID} devgroup \
    && adduser --system --uid ${USER_UID} devuser

USER devuser

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
