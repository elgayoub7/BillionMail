#!/usr/bin/env python3
"""Patch batch_mail.go API structs"""
import os

path = os.path.join(os.path.dirname(__file__), "api", "batch_mail", "v1", "batch_mail.go")

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

if "SendDelay" in content:
    print("ALREADY PATCHED")
    exit(0)

schedule_fields_create = """
\tSendDelay           int    `json:"send_delay"            dc:"Delay between emails in seconds" default:"0"`
\tScheduleStartHour   int    `json:"schedule_start_hour"   dc:"Schedule start hour 0-23" default:"0"`
\tScheduleEndHour     int    `json:"schedule_end_hour"     dc:"Schedule end hour 0-24" default:"24"`
\tScheduleDays        string `json:"schedule_days"         dc:"Allowed weekdays JSON [1-7] ISO" default:"[1,2,3,4,5,6,7]"`
\tSenderPool          string `json:"sender_pool"           dc:"Sender pool JSON for rotation" default:"[]"`
\tDailyLimitPerSender int    `json:"daily_limit_per_sender" dc:"Max emails per sender per day" default:"0"`
"""

schedule_fields_update = """
\tSendDelay           int    `json:"send_delay"            dc:"Delay between emails in seconds"`
\tScheduleStartHour   int    `json:"schedule_start_hour"   dc:"Schedule start hour 0-23"`
\tScheduleEndHour     int    `json:"schedule_end_hour"     dc:"Schedule end hour 0-24"`
\tScheduleDays        string `json:"schedule_days"         dc:"Allowed weekdays JSON [1-7]"`
\tSenderPool          string `json:"sender_pool"           dc:"Sender pool JSON for rotation"`
\tDailyLimitPerSender int    `json:"daily_limit_per_sender" dc:"Max emails per sender per day"`
"""

# Patch CreateTaskReq (line 157-158)
old1 = '\tTagLogic string `json:"tag_logic" v:"in:AND,OR,NOT" dc:"tag logic (AND: must have all tags, OR: have any tag, NOT)" default:"AND"`\n}'
new1 = '\tTagLogic string `json:"tag_logic" v:"in:AND,OR,NOT" dc:"tag logic (AND: must have all tags, OR: have any tag, NOT)" default:"AND"`' + schedule_fields_create + '}'

if old1 in content:
    content = content.replace(old1, new1, 1)
    print("OK: Patched CreateTaskReq")
else:
    print("ERROR: CreateTaskReq pattern not found")

# Patch UpdateTaskInfoReq (line 283-284)
old2 = '\tTagLogic      string `json:"tag_logic" v:"in:AND,OR,NOT" dc:"tag logic (AND: must have all tags, OR: have any tag, NOT)"`\n}'
new2 = '\tTagLogic      string `json:"tag_logic" v:"in:AND,OR,NOT" dc:"tag logic (AND: must have all tags, OR: have any tag, NOT)"`' + schedule_fields_update + '}'

if old2 in content:
    content = content.replace(old2, new2, 1)
    print("OK: Patched UpdateTaskInfoReq")
else:
    print("ERROR: UpdateTaskInfoReq pattern not found")

with open(path, "w", encoding="utf-8", newline="\n") as f:
    f.write(content)
