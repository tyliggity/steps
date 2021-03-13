#!/bin/bash
set -eo pipefail

if [[ -z "${SP_PREVIOUS_STEP_OUTPUT_URL}" ]]; then
	echo "Missing \"SP_PREVIOUS_STEP_OUTPUT_URL\" environment variable" >/dev/stderr
	exit 1
fi

[[ -z $SP_DEBUG  ]] || echo "SP_PREVIOUS_STEP_OUTPUT_URL="${SP_PREVIOUS_STEP_OUTPUT_URL} > /dev/stderr

if [[ ${SP_PREVIOUS_STEP_OUTPUT_URL} = "gs://"* ]]; then
    for i in {1..5}
    do
	    /usr/local/bin/sp-base-step cat -url $SP_PREVIOUS_STEP_OUTPUT_URL && break || sleep 5
    done
else
    echo "Unsupported log url prefix provided"
    exit 1
fi