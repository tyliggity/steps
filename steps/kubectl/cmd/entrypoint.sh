#!/bin/sh

# This file is here because there is a problem with spawning /bin/sh and passing all the args (it's truncating the first argument).
set -eo pipefail

/cmd $@ 2>&1 | sp-base-step format