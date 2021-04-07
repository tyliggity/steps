FROM debian:bullseye-slim

ARG BASEDIR=/go/src/github.com/stackpulse

ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/ .

RUN apt-get update
RUN apt-get install -y python3 pip gcc python-dev libkrb5-dev libssl-dev
RUN pip install pywinrm[kerberos]
RUN pip install pywinrm[credssp]
RUN chmod +x entrypoint.py

ENTRYPOINT ["./entrypoint.py"]
