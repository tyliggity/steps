# Globals
ARG BASEDIR=/go/src/github.com/stackpulse
ARG BUILD_BRANCH
ARG STEP_BASEPATH

ARG BUILD_BRANCH
FROM us-docker.pkg.dev/stackpulse/public/base:$BUILD_BRANCH AS stepsBase
FROM golang:1.14-alpine

COPY --from=stepsBase /usr/local/bin/ /usr/local/bin/

RUN apk add --no-cache curl
RUN curl -LO https://coralogix-public.s3-eu-west-1.amazonaws.com/scanner/v0.2.1/scanner-Linux.gz
RUN gunzip scanner-Linux.gz
RUN chmod +x scanner-Linux
ENTRYPOINT [ "scanner-Linux" ]