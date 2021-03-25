#! /usr/bin/env python

"""
    Used to build a step
"""

import helpers
import sys
import os


def step_family_path(split_step_path):
    family_dir = os.path.join(os.getcwd(), "/".join([".."] * (len(split_step_path) - 1)))
    return os.path.normpath(family_dir)


def main():
    step_rel_path = helpers.get_step_rel_path(os.getcwd())
    split_step_path = [x for x in step_rel_path.split("/") if x != ""]
    image_tag = helpers.get_step_docker_image_tag(step_rel_path)

    if not helpers.docker_build(
            image_tag,
        "Dockerfile",
        step_family_path(split_step_path),
        ["--build-arg", "BASE_BRANCH=latest",
         "--build-arg", "CURRENT_BRANCH=" + helpers.get_current_branch(),
         "--build-arg", "STEP_BASEPATH=" + "/".join(split_step_path[1:])]):
        return 1

    if not helpers.docker_push(image_tag):
        return 1


if __name__ == "__main__":
    sys.exit(main())
