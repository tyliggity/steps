FROM confluentinc/cp-kafka:5.0.1
ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/delete.sh /delete.sh
RUN chmod +x /delete.sh
ENTRYPOINT ["bash", "-c", "/delete.sh"]