if [[ -z "${TOPIC_NAME}" ]] ; then
  topic_arg=""
else
  topic_arg="--topic ${TOPIC_NAME}"
fi
/usr/bin/kafka-topics --zookeeper ${ZOOKEPER_HOST} --describe ${topic_arg}