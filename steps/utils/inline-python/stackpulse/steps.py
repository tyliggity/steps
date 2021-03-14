import subprocess
PREV_OUTPUT_CACHE = None
def load_previous_step_output():
    global PREV_OUTPUT_CACHE
    if PREV_OUTPUT_CACHE is None:
        PREV_OUTPUT_CACHE = subprocess.run(['cat-previous-step-output'], stdout=subprocess.PIPE).stdout
    return PREV_OUTPUT_CACHE
