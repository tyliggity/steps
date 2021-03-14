FROM ubuntu
RUN apt-get update && apt-get install -y \
  python \
  redis-server \
  black
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/hset.py /tmp/

RUN black --check /tmp

ENTRYPOINT ["python", "/tmp/hset.py"]