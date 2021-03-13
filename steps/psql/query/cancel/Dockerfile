# Globals
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BINARY_NAME=query-cancel
ARG BUILD_BRANCH
ARG STEP_BASEPATH

# 1. Base
FROM us-docker.pkg.dev/stackpulse/public/base:$BUILD_BRANCH AS base

# 2. Builder
FROM golang:1.14-alpine as builder
ARG BUILD_BRANCH
ARG STEP_BASEPATH
ARG BINARY_NAME
ARG BASEDIR

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN GOOS=linux GOARCH=amd64 go build -o /${BINARY_NAME}

# 3. Production
FROM postgres:12-alpine
ARG BINARY_NAME
RUN apk --update add openssh-client
COPY --from=builder /${BINARY_NAME} .
RUN chmod +x ${BINARY_NAME}
COPY --from=base /usr/local/bin/ /usr/local/bin/

ENV STEP_NAME=${BINARY_NAME}
ENTRYPOINT ["/bin/sh", "-c", "set -eo pipefail; /${STEP_NAME} 2>&1 | sp-base-step format"]