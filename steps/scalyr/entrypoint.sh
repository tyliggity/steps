#!/bin/sh

scalyr "$@" \
--token=$SCALYR_TOKEN  \
--output=json-pretty
