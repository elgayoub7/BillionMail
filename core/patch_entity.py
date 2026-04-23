#!/usr/bin/env python3
"""Patch batch_mail.go entity to add cold mail scheduling fields"""
import os, re

path = os.path.join(os.path.dirname(__file__), "internal", "model", "entity", "batch_mail.go")

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

if "SendDelay" in content:
    print("ALREADY PATCHED")
    exit(0)

# Find the closing brace of EmailTask struct after UseTagFilter
# Pattern: UseTagFilter line, then closing brace
pattern = r'(\tUseTagFilter\s+int\s+`json:"use_tag_filter"[^`]*`)\n\}'
match = re.search(pattern, content)
if not match:
    print("ERROR: UseTagFilter pattern not found")
    exit(1)

insert_after = match.group(0)
new_fields = insert_after + """

\t// Cold mail scheduling fields
\tSendDelay           int    `json:"send_delay"            dc:"Delay between emails in seconds (0=no delay)" orm:"send_delay"`
\tScheduleStartHour   int    `json:"schedule_start_hour"   dc:"Allowed sending start hour 0-23 (0=no restriction)" orm:"schedule_start_hour"`
\tScheduleEndHour     int    `json:"schedule_end_hour"     dc:"Allowed sending end hour 0-24 (24=no restriction)" orm:"schedule_end_hour"`
\tScheduleDays        string `json:"schedule_days"         dc:"Allowed weekdays ISO [1-7], 1=Mon, 7=Sun" orm:"schedule_days"`
\tSenderPool          string `json:"sender_pool"           dc:"JSON array of {email,name} for rotation" orm:"sender_pool"`
\tDailyLimitPerSender int    `json:"daily_limit_per_sender" dc:"Max emails per sender per day (0=unlimited)" orm:"daily_limit_per_sender"`
\tCurrentSenderIndex  int    `json:"current_sender_index"  dc:"Current sender index for round-robin" orm:"current_sender_index"`
}"""

content = content[:match.start()] + new_fields + content[match.end():]
print("OK: Added cold mail scheduling fields to EmailTask")

# Patch AfterFind - find the function and add defaults before the closing brace
af_pattern = r'(func \(e \*EmailTask\) AfterFind\(\) \{\n\tif e\.TagIdsRaw != "" \{\n\t\tvar tagIds \[\]int\n\t\terr := json\.Unmarshal\(\[\]byte\(e\.TagIdsRaw\), &tagIds\)\n\t\tif err == nil \{\n\t\t\te\.TagIds = tagIds\n\t\t\}\n\t\}\n)(})'
af_match = re.search(af_pattern, content)
if af_match:
    new_af = af_match.group(1) + """\tif e.ScheduleStartHour == 0 && e.ScheduleEndHour == 0 {
\t\te.ScheduleEndHour = 24
\t}
\tif e.ScheduleDays == "" {
\t\te.ScheduleDays = "[1,2,3,4,5,6,7]"
\t}
\tif e.SenderPool == "" {
\t\te.SenderPool = "[]"
\t}
}"""
    content = content[:af_match.start()] + new_af + content[af_match.end():]
    print("OK: Updated AfterFind with scheduling defaults")
else:
    print("WARN: Could not patch AfterFind")

with open(path, "w", encoding="utf-8", newline="\n") as f:
    f.write(content)
