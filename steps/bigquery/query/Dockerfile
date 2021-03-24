ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/src/stackpulse
FROM us-docker.pkg.dev/stackpulse/public/base:$BASE_BRANCH AS base

FROM golang:1.15-alpine as builder
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/src/stackpulse
WORKDIR $BASEDIR
COPY . $BASEDIR
WORKDIR $BASEDIR/$STEP_BASEPATH
RUN GOOS=linux GOARCH=amd64 go build -o /bigquery-query

FROM alpine
ARG BASE_BRANCH
ARG STEP_BASEPATH
ARG BASEDIR=/src/stackpulse
COPY --from=builder /bigquery-query .
RUN chmod +x bigquery-query
COPY --from=base /usr/local/bin/ /usr/local/bin/
COPY $STEP_BASEPATH/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
