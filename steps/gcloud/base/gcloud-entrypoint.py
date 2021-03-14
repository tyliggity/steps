import subprocess
import os
import base64

auth_code = os.getenv("AUTH_CODE")
command = os.getenv("COMMAND")
project_id = os.getenv("PROJECT_ID")

raw_account_serivce_id = base64.decodestring(auth_code)
f = open("/tmp/auth.json", "w")
f.write(raw_account_serivce_id + "\n")
f.close()

res = subprocess.check_output(
    ["gcloud", "auth", "activate-service-account", "--key-file", "/tmp/auth.json"]
)
for line in res.splitlines():
    print(line)

res = subprocess.check_output(["gcloud", "config", "set", "project", project_id])
for line in res.splitlines():
    print(line)

splittedcommand = command.split(" ")
res = subprocess.check_output(splittedcommand)
for line in res.splitlines():
    print(line)
