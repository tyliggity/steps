FROM python:3 as build-env

WORKDIR /src/
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH} /src/

RUN pip install black
RUN black --check .

FROM google/cloud-sdk
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/gcloud-entrypoint.py /tmp/
ENTRYPOINT ["python", "/tmp/gcloud-entrypoint.py"]
