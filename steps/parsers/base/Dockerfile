FROM google/cloud-sdk:alpine
RUN apk add --update \
curl \
bash \
python3 
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/cat-previous-step-output.sh /bin/cat-previous-step-output
COPY ${STEP_BASEPATH}/base-funcs.sh /import/base
