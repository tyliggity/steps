import os

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


# get steps with Dockerfile
def get_steps():
    for root, dirs, files in os.walk(BASE_DIR, topdown=False):
        skip_dir = False
        for id in IGNORE_DIRS:
            if root.startswith(id):
                skip_dir = True
                break

        if skip_dir:
            continue

        has_dockerfile = False
        for f in files:
            if f == "Dockerfile":
                has_dockerfile = True

        if not has_dockerfile:
            continue

        yield root