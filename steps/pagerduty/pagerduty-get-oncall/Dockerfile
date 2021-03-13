FROM us-docker.pkg.dev/stackpulse/public/0_base AS stepsBase
FROM debian:jessie-slim
RUN apt-get update && apt-get install -y ca-certificates curl
COPY --from=stepsBase /usr/local/bin/ /usr/local/bin/
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/bin/sh", "-c", "/entrypoint.sh 2>&1 | sp-base-step format"]
