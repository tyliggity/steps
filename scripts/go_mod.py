#! /usr/bin/env python
"""
    Runs
        - go get ./...
        - go mod download
        - go mod tidy
        in every go step directory
"""

import os
import sys

IGNORE_DIRS = ["./.git", "./.idea"]
GOMOD = "go.mod"
BASE_DIR = "./steps"

# return path to all go steps
def get_go_steps():
    for root, dirs, files in os.walk(BASE_DIR, topdown=False):
        skip_dir = False
        for id in IGNORE_DIRS:
            if root.startswith(id):
                skip_dir = True
                break

        if skip_dir:
            continue

        has_go_mod = False
        for f in files:
            if f == GOMOD:
                has_go_mod = True

        if not has_go_mod:
            continue

        yield root


def main():
    for step in get_go_steps():
        cur_dir = os.getcwd()
        print("> Running in: %s" % (step,))
        os.chdir(step)
        os.system("go get ./...")
        os.system("go mod download")
        os.system("go mod tidy")
        print("Done")
        os.chdir(cur_dir)

    return 0

if __name__ == "__main__":
    sys.exit(main())