import constants
import os
import subprocess


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
def get_step_docker_image(step_path):
    return os.path.join(constants.CONTAINER_REGISTRY, step_path)


# Get docker image tag
def get_step_docker_image_tag(step_path):
    tag = get_current_branch()
    if tag == "master":
        tag = "latest"

    return get_step_docker_image(step_path) + ":" + tag


# Build step docker image
def docker_build(tag, dockerfile="Dockerfile", root_path=".", args=[]):
    print("Building docker image %s" % (tag,))
    cmd = ["docker", "build"]
    cmd += args
    cmd += ["--iidfile", constants.CONTAINER_ID_FILE,
           "-t", tag, "-f", dockerfile, root_path]

    print("Running: %s" " ".join(cmd))
    if 0 != os.system(" ".join(cmd)):
        raise Exception("build command failed")

    return True


def docker_push(image):
    if os.getenv("CI") != "true":
        print("Skipping push for local build")
        return True

    print("Pushing docker image %s" % (image, ))
    if 0 != os.system("docker push %s" % (image, )):
        raise Exception("docker push failed")
    return True

