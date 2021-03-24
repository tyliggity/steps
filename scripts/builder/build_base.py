#! /usr/bin/env python
"""
    This script is used to build ./base image
"""

import helpers
import sys

STEP_PATH = "base"

def main():
    image_tag = helpers.get_step_docker_image_tag(STEP_PATH)

    if not helpers.docker_build(image_tag):
        return 1

    if not helpers.docker_push(image_tag):
        return 1


if __name__ == "__main__":
    sys.exit(main())