#!/bin/bash

export END_MARKER="<-- END -->"

function error {
	(echo -n "[-] "; echo $*) > /dev/stderr
	exit 1
}

function debug {
	[[ -z "${DEBUG}" ]] || (echo -n "[*] "; echo $*) > /dev/stderr
}

function assert-environment-variable {
	[[ -z "$1" ]] && error "Missing required environment variable \"$2\"."
}

function get-input {
	if [[ -z "${INPUT}" ]] ; then
		debug "Using previous step output as tranformation input.";
		cat-previous-step-output;
	else
		debug "Using \"INPUT\" environment variable as transformation input.";
		echo "$INPUT";
	fi
}