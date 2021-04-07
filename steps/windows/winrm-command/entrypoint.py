#!/usr/bin/env python3

import env
import shlex

END_STRING = "<-- END -->"

HOST = os.environ["HOST"]
USERNAME = os.environ["USERNAME"]
PASSWORD = os.environ["PASSOWRD"]
COMMAND = os.environ["COMMAND"]

commands = shlex.split(COMMAND)
s = winrm.Session(HOST, auth=(USERNAME, PASSWORD))
r = s.run_cmd(commands[0], commands[1:])

print(r.status_code)
print(r.std_out)
print(r)
output_dict = {"output": r.std_out, "status": r.status_code}
print(f"{END_STRING}{json.dumps(output_dict)}")
