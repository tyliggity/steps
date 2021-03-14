FROM confluentinc/cp-kafka:5.0.1
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/describe.sh /describe.sh
RUN chmod +x /describe.sh
ENTRYPOINT ["bash", "-c", "/describe.sh"]