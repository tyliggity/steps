FROM ubuntu
RUN apt-get update && apt-get install -y \
  python \
  redis-server \
  black
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/big-keys.py /tmp/

RUN black --check /tmp

ENTRYPOINT ["python", "/tmp/big-keys.py"]