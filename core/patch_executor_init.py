#!/usr/bin/env python3
"""Fix: add scheduleChecker init to NewTaskExecutor"""
import os

path = os.path.join(os.path.dirname(__file__), "internal", "service", "batch_mail", "task_executor.go")

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

if "scheduleChecker:" in content:
    print("ALREADY PATCHED")
    exit(0)

old = """\t\trateController: NewSimpleRateController(1000),
\t}

\treturn executor
}"""

new = """\t\trateController:  NewSimpleRateController(1000),
\t\tscheduleChecker: NewScheduleChecker(),
\t}

\treturn executor
}"""

if old in content:
    content = content.replace(old, new, 1)
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(content)
    print("OK: scheduleChecker initialized in NewTaskExecutor")
else:
    print("ERROR: pattern not found")
    # Debug
    idx = content.find("rateController: NewSimpleRateController")
    if idx >= 0:
        print(repr(content[idx-50:idx+150]))
