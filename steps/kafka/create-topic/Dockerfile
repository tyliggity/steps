FROM confluentinc/cp-kafka:5.0.1
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/create.sh /create.sh
RUN chmod +x /create.sh
ENTRYPOINT ["bash", "-c", "/create.sh"]