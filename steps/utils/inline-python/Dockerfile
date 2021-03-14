FROM us-docker.pkg.dev/stackpulse/public/parsers/base:SP-1815

WORKDIR /usr/src/app
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/ .

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
