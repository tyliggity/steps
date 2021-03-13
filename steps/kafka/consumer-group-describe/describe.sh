if [[ -z "${MEMBERS}" ]] ; then
  members_arg=""
else
  members_arg="--members"
fi

if [[ -z "${STATE}" ]] ; then
  state_arg=""
else
  state_arg="--state"
fi

/usr/bin/kafka-consumer-groups --bootstrap-server ${BOOTSTRAP_HOST} --group ${GROUP_NAME} --describe ${members_arg} ${state_arg}