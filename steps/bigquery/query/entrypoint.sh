#!/bin/sh

set -eo pipefail

readonly SA_FILE=/var/gcp.json

function check_auth() {
	if [ ! -z $AUTH_CODE ]
	then
		echo $AUTH_CODE | base64 -d > $SA_FILE
		export GOOGLE_APPLICATION_CREDENTIALS=$SA_FILE
	fi
}

function main() {
	check_auth
	/bigquery-query 2>&1 | sp-base-step format
}

# execute here.
main
