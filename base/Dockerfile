FROM golang:1.14 as builder

ARG BASE_STEP_REPO=/go/src/github.com/stackpulse/base

WORKDIR ${BASE_STEP_REPO}
COPY . .

WORKDIR ${BASE_STEP_REPO}
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /sp-base-step .

FROM scratch
ARG BASE_STEP_REPO=/go/src/github.com/stackpulse/base
WORKDIR /usr/local/bin
COPY --from=builder /sp-base-step .
COPY --from=builder ${BASE_STEP_REPO}/bash/ .
