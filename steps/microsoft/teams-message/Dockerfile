FROM debian:jessie-slim
RUN apt-get update && apt-get install -y ca-certificates curl
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/entrypoint.sh /
RUN chmod +x /entrypoint.sh

CMD /entrypoint.sh
