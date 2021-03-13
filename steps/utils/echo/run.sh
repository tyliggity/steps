#!/bin/sh

if [ -z "$MESSAGE" ]; then
	/bin/echo $@;
else
	/bin/echo $MESSAGE;
fi
