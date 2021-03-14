FROM python:3-alpine

# Install curl
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

# Install bash for a simple entry point
RUN apk add --no-cache --upgrade bash

# Install Scalyr command-line tool (implemented in Python)
RUN curl https://raw.githubusercontent.com/scalyr/scalyr-tool/master/scalyr > /usr/local/bin/scalyr
RUN chmod +x /usr/local/bin/scalyr

# Copy and run the entry point shell script
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/bin/sh","/entrypoint.sh"]
