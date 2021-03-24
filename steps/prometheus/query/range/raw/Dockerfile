# Globals
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BINARY_NAME=prometheus-range-raw
ARG BASE_BRANCH
ARG STEP_BASEPATH

# 1. Base
FROM us-docker.pkg.dev/stackpulse/public/base:$BASE_BRANCH AS base

# 2. Builder
FROM golang:1.14-alpine as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BINARY_NAME
ARG BASEDIR

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN GOOS=linux GOARCH=amd64 go build -o /${BINARY_NAME}

# 3. Production
FROM alpine:3.7
RUN apk add --update libc6-compat
ARG BINARY_NAME
COPY --from=builder /${BINARY_NAME} .
RUN chmod +x ${BINARY_NAME}
COPY --from=base /usr/local/bin/ /usr/local/bin/

ENV STEP_NAME=${BINARY_NAME}
ENTRYPOINT ["/bin/sh", "-c", "/${STEP_NAME} 2>&1 | sp-base-step format"]