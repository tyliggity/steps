#!/usr/bin/env python3

s = winrm.Session('windows-host.example.com', auth=('john.smith', 'secret'))
r = s.run_cmd('ipconfig', ['/all'])
print(r.status_code)
print(r.std_out)
