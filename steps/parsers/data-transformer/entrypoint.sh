#!/bin/bash

source /import/base

assert-environment-variable "${COMMAND}" COMMAND


COMMAND="get-input | ${COMMAND}"

debug "About to run: ${COMMAND}"

RESPONSE=`eval "${COMMAND}"`

debug "Response = ${RESPONSE}"
if [ -z "$RESPONSE" ]; then
    error "Empty response"
    exit 1
fi


if [ -z "$DEBUG" ]; then
    JSON=`jq -n --argjson response "$RESPONSE" '{output: $response }' 2>/dev/null`
else
    JSON=`jq -n --argjson response "$RESPONSE" '{output: $response }'`
fi

if [ "$?" -ne 0 ]; then
    JSON=`jq -n --arg response "$RESPONSE" '{output: $response }'`
fi 



echo "${END_MARKER}${JSON}"
