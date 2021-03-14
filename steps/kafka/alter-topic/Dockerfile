FROM confluentinc/cp-kafka:5.0.1
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/alter.sh /alter.sh
RUN chmod +x /alter.sh
ENTRYPOINT ["bash", "-c", "/alter.sh"]