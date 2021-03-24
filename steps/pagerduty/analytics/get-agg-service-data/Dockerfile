ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse

FROM us-docker.pkg.dev/stackpulse/public/base:$BASE_BRANCH AS base
FROM golang:1.15-alpine as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN [ -z `gofmt -l .` ]
RUN GOOS=linux GOARCH=amd64 go build -o /pagerduty-analytics-get-agg-service-data


FROM alpine:3.7
WORKDIR /opt/stackpulse
COPY --from=builder /pagerduty-analytics-get-agg-service-data /opt/stackpulse/
RUN chmod +x pagerduty-analytics-get-agg-service-data
COPY --from=base /usr/local/bin/ /usr/local/bin/

ENTRYPOINT ["/bin/sh", "-c", "./pagerduty-analytics-get-agg-service-data 2>&1 | sp-base-step format"]