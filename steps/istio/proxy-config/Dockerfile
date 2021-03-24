# 1. Globals
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BINARY_NAME=pc-wrapper

# 2. Builder
FROM us-docker.pkg.dev/stackpulse/public/base:$BASE_BRANCH AS base
FROM golang:1.14-buster as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BINARY_NAME
ARG BASEDIR=/go/src/github.com/stackpulse

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN GOOS=linux GOARCH=amd64 go build -o /${BINARY_NAME} && chmod +x /${BINARY_NAME}

# 3. Production
FROM istio/istioctl:1.6.10
ARG BINARY_NAME
ENV STEP_NAME=${BINARY_NAME}
COPY --from=builder /${BINARY_NAME} .
COPY --from=base /usr/local/bin/ /usr/local/bin/
ENTRYPOINT ["/bin/sh", "-c", "set -eo pipefail; /${STEP_NAME} 2>&1 | sp-base-step format"]