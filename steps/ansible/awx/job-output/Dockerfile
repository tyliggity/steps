# 1. Globals
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BINARY_NAME=awx-job-output

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
FROM python:3.9.1-alpine
RUN pip install "https://github.com/ansible/awx/archive/17.0.1.tar.gz#egg=awxkit&subdirectory=awxkit"
ARG BINARY_NAME
ENV STEP_NAME=${BINARY_NAME}
COPY --from=builder /${BINARY_NAME} .
COPY --from=base /usr/local/bin/ /usr/local/bin/
ENTRYPOINT ["/bin/sh", "-c", "set -eo pipefail; /${STEP_NAME} 2>&1 | sp-base-step format"]