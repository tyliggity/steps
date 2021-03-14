FROM confluentinc/cp-kafka:5.0.1
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/list.sh /list.sh
RUN chmod +x /list.sh
ENTRYPOINT ["bash", "-c", "/list.sh"]