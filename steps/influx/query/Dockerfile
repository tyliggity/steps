ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse

FROM us-docker.pkg.dev/stackpulse/public/base:$BASE_BRANCH AS base
FROM golang:1.14-alpine as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/go/src/github.com/stackpulse

WORKDIR ${BASEDIR}/
COPY . .
WORKDIR ${BASEDIR}/${STEP_BASEPATH}

RUN GOOS=linux GOARCH=amd64 go build -o /influx-query


FROM influxdb:1.7.10-alpine
COPY --from=builder /influx-query .
RUN chmod +x influx-query
COPY --from=base /usr/local/bin/ /usr/local/bin/

ENTRYPOINT ["/bin/sh", "-c", "/influx-query 2>&1 | sp-base-step format"]