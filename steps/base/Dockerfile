FROM golang:1.14 as builder

ARG BASE_STEP_REPO=/go/src/github.com/stackpulse/

WORKDIR ${BASE_STEP_REPO}
COPY . .

WORKDIR ${BASE_STEP_REPO}/steps/base
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /sp-base-step .

FROM scratch
ARG BASE_STEP_REPO=/go/src/github.com/stackpulse/steps/base
WORKDIR /usr/local/bin
COPY --from=builder /sp-base-step .
COPY --from=builder ${BASE_STEP_REPO}/bash/ .
