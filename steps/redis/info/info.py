#!/usr/bin/env python

import subprocess
import os
import base64
import json

redis_ip = os.getenv("REDIS_IP")
redis_password = os.getenv("REDIS_PASSWORD")
redis_url = os.getenv("REDIS_URL")

cmd = ["redis-cli", "--no-auth-warning"]
if redis_ip != None:
    cmd.extend(["-h", redis_ip])
    if redis_password != None:
        cmd.extend(["-a", redis_password])
if redis_url != None:
    cmd.extend(["-u", redis_url])

cmd.append("info")

formattedStepResponse = {}
res = subprocess.check_output(cmd)

final_object = {}
splittedLines = res.splitlines()


def build_sector(sector_name):
    sector_index = 1
    final_object[sector_name] = {}
    while (
        len(splittedLines) > lineIndex + sector_index
        and len(splittedLines[lineIndex + sector_index]) != 0
    ):
        final_object[sector_name][
            splittedLines[lineIndex + sector_index].split(":")[0]
        ] = splittedLines[lineIndex + sector_index].split(":")[1]
        sector_index = sector_index + 1

    return sector_index


lineIndex = 0
while lineIndex < len(splittedLines):
    if splittedLines[lineIndex][0] == "#":
        lineIndex = lineIndex + build_sector(splittedLines[lineIndex].split()[1])
    lineIndex = lineIndex + 1

result = json.dumps({"info": final_object, "multiline": res})

print("<-- END -->")
print(result)
