FROM debian:jessie-slim
RUN apt-get update && apt-get install -y ca-certificates curl openssh-client jq unzip

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install


ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/entrypoint.sh /
RUN chmod +x /entrypoint.sh

CMD /entrypoint.sh
