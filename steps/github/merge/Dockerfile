FROM us-docker.pkg.dev/stackpulse/public/0_base AS stepsBase
FROM alpine:3.5
RUN apk update
RUN apk add git
RUN apk add openssh

COPY --from=stepsBase /usr/local/bin/ /usr/local/bin/
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/bin/sh", "-c", "/entrypoint.sh"]
