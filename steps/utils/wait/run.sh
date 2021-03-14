#!/bin/sh

if [ -z "$SECONDS" ]; then
	/bin/sleep 10;
else
	/bin/sleep $SECONDS;
fi
