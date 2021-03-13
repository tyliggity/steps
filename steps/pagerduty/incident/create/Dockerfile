ARG BUILD_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse

FROM us-docker.pkg.dev/stackpulse/public/base:$BUILD_BRANCH AS base
FROM golang:1.14-alpine as builder
ARG BUILD_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN [ -z `gofmt -l .` ]
RUN GOOS=linux GOARCH=amd64 go build -o /pagerduty-incident-create


FROM alpine:3.7
WORKDIR /opt/stackpulse
COPY --from=builder /pagerduty-incident-create /opt/stackpulse/
RUN chmod +x pagerduty-incident-create
COPY --from=base /usr/local/bin/ /usr/local/bin/

ENTRYPOINT ["/bin/sh", "-c", "./pagerduty-incident-create 2>&1 | sp-base-step format"]