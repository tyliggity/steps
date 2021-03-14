#!/usr/bin/env python

import subprocess
import os
import base64
import json

redis_ip = os.getenv("REDIS_IP")
redis_password = os.getenv("REDIS_PASSWORD")
redis_url = os.getenv("REDIS_URL")

key = os.getenv("KEY")
field = os.getenv("FIELD")
value = os.getenv("VALUE")

cmd = ["redis-cli", "--no-auth-warning"]
if redis_ip != None:
    cmd.extend(["-h", redis_ip])
    if redis_password != None:
        cmd.extend(["-a", redis_password])
if redis_url != None:
    cmd.extend(["-u", redis_url])

cmd.extend(["hset", key, field, value])

res = subprocess.check_output(cmd)

for line in res.splitlines():
    print(line)
