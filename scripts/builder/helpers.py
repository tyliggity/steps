import constants
import os
import subprocess
import yaml


def system(cmd):
    print("> Running: %s" % (cmd,))
    return os.system(cmd)

# Get current branch using git cli tool
def get_current_branch():
    return subprocess.check_output(["git", "rev-parse", "--abbrev-ref", "HEAD"]).strip().decode("utf-8")


# Get step relative path (to /steps/)
def get_step_rel_path(path):
    s = path.split(constants.STEPS_ROOT)
    if len(s) != 2:
        raise Exception("cant find delimiter in step path '%s'" % (path,))

    return s[1]


# Get step docker image path
def get_step_docker_repository(step_path):
    return os.path.join(constants.CONTAINER_REGISTRY, step_path)


# Get manifest version
def get_manifest_version(step_path):
    manifest_path = os.path.join(step_path, constants.MANIFEST_FILENAME)
    try:
        f = open()
        y = yaml.safe_load(f)
        f.close()
        version = y['metadata']['version']
    except:
        print("[X] Failed parsing version from %s" % (manifest_path,))
        raise Exception("failed parsing manifest")

    if version == "":
        raise Exception("no valid version set for step: %s" % (manifest_path,))


# Build full docker image path given a repo and tag
def docker_image_tag(repo, tag):
    return repo + ":" + tag


# return the relevant tags for the given branch
# on side branch we tag with current branch on master we tag with latest and imageVersion
def get_step_image_tags(step_path):
    current_branch = get_current_branch()
    if current_branch == "master":
        return ["latest", get_manifest_version(step_path)]
    else:
        return [current_branch]


# Build step docker image
def docker_build(tag, dockerfile="Dockerfile", root_path=".", args=[]):
    print("Building docker image %s" % (tag,))
    cmd = ["docker", "build"]
    cmd += args
    cmd += ["--iidfile", constants.CONTAINER_ID_FILE,
            "-t", tag, "-f", dockerfile, root_path]

    if 0 != system(" ".join(cmd)):
        raise Exception("build command failed")

    return True


# Tags the given docker image with the tags
def docker_tag(docker_repo, current_tag, new_tag):
    old_tag = docker_image_tag(docker_repo, current_tag)
    new_image_with_tag = docker_image_tag(docker_repo, new_tag)
    print("> Tagging %s -> %s" % (old_tag, new_tag))
    if 0 != system("docker tag %s %s" % (old_tag, new_image_with_tag)):
        raise Exception("docker tag failed")

    return True


def docker_push(image):
    if os.getenv("CI") != "true":
        print("Skipping push for local build")
        return True

    print("Pushing docker image %s" % (image,))
    if 0 != system("docker push %s" % (image,)):
        raise Exception("docker push failed")
    return True
