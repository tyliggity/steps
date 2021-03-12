import subprocess
import os
import base64
import json

auth_code   = os.getenv('AUTH_CODE')
project_id  = os.getenv('PROJECT_ID')
gcloud_command = os.getenv('COMMAND')


# We expect to receive the key file JSON in a Basc64-encoded environment variable.
# Lets decode it and save it to local file for use with gcloud cli
raw_account_serivce_id = base64.decodestring(auth_code+"===")
f = open("/tmp/auth.json", "w")
f.write(raw_account_serivce_id + "\n")
f.close()

# Try to authenticate with a key file and print the output to the stdout
res = subprocess.check_output(["gcloud", "auth", "activate-service-account", "--key-file", "/tmp/auth.json"])
for line in res.splitlines():
    print line

# Set the the current project
res = subprocess.check_output(["gcloud", "config", "set", "project", project_id])
for line in res.splitlines():
    print line


splittedcommand = gcloud_command.split(" ")

print splittedcommand
print "Splitting Quotes"

# Process and format step output
stepoutput = "<-- END -->{\"output\": \""
res = subprocess.check_output(splittedcommand)

for line in res.splitlines():
    stepoutput = stepoutput + json.dumps(line).strip("\"")
    stepoutput = stepoutput + "\\n"
    print line
    
stepoutput = stepoutput + "\"}"

print stepoutput

