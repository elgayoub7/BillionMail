#!/usr/bin/env python3
"""Patch task_executor.go to add scheduling, delay, sender rotation, daily limits"""
import os, re

path = os.path.join(os.path.dirname(__file__), "internal", "service", "batch_mail", "task_executor.go")

with open(path, "r", encoding="utf-8") as f:
    content = f.read()

changes = 0

# 1. Add senderRotation and scheduleChecker fields to TaskExecutor struct
old_struct = """\t// circuit breaker
\tconsecutiveFailures atomic.Int64

\t// pause/resume control
\tpauseChan  chan struct{}
\tresumeChan chan struct{}
}"""

new_struct = """\t// circuit breaker
\tconsecutiveFailures atomic.Int64

\t// pause/resume control
\tpauseChan  chan struct{}
\tresumeChan chan struct{}

\t// cold mail features
\tsenderRotation *SenderRotation
\tscheduleChecker *ScheduleChecker
}"""

if old_struct in content and "senderRotation" not in content:
    content = content.replace(old_struct, new_struct, 1)
    changes += 1
    print("OK: Added senderRotation + scheduleChecker to TaskExecutor struct")
else:
    print("SKIP: TaskExecutor struct (already patched or not found)")

# 2. Initialize scheduleChecker in NewTaskExecutor
old_init = """\t\tstartTime:      time.Now(),
\t\tpauseChan:      make(chan struct{}, 1),
\t\tresumeChan:     make(chan struct{}, 1),
\t\trateController: NewSimpleRateController(1000),
\t}"""

new_init = """\t\tstartTime:       time.Now(),
\t\tpauseChan:       make(chan struct{}, 1),
\t\tresumeChan:      make(chan struct{}, 1),
\t\trateController:  NewSimpleRateController(1000),
\t\tscheduleChecker: NewScheduleChecker(),
\t}"""

if old_init in content and "scheduleChecker" not in content:
    content = content.replace(old_init, new_init, 1)
    changes += 1
    print("OK: Initialized scheduleChecker in NewTaskExecutor")
else:
    print("SKIP: NewTaskExecutor init")

# 3. Initialize senderRotation in ProcessTask, after getting task info
old_process = """\t// configure rate controller
\te.configureRateController(task)"""

new_process = """\t// initialize sender rotation (if sender pool configured)
\tif task.SenderPool != "" && task.SenderPool != "[]" {
\t\te.senderRotation = NewSenderRotation(task.SenderPool, task.DailyLimitPerSender, task.Id, task.CurrentSenderIndex)
\t\tg.Log().Infof(ctx, "task %d: sender rotation enabled, pool size: %d, daily limit: %d", task.Id, e.senderRotation.PoolSize(), task.DailyLimitPerSender)
\t}

\t// configure rate controller
\te.configureRateController(task)"""

if old_process in content and "initialize sender rotation" not in content:
    content = content.replace(old_process, new_process, 1)
    changes += 1
    print("OK: Added sender rotation init in ProcessTask")
else:
    print("SKIP: ProcessTask sender rotation init")

# 4. Add schedule check in processTaskRecipients main loop (before getting batch)
old_loop = """\tfor {
\t\t// check if context is canceled
\t\tselect {
\t\tcase <-ctx.Done():
\t\t\tg.Log().Info(ctx, "context canceled, stop task execution:", ctx.Err())
\t\t\treturn ctx.Err()
\t\tdefault:
\t\t}

\t\t// check pause status"""

new_loop = """\tfor {
\t\tselect {
\t\tcase <-ctx.Done():
\t\t\tg.Log().Info(ctx, "context canceled, stop task execution:", ctx.Err())
\t\t\treturn ctx.Err()
\t\tdefault:
\t\t}

\t\t// schedule check: time window + day of week
\t\tif !e.scheduleChecker.IsWithinSchedule(task) {
\t\t\tif !e.scheduleChecker.SleepUntilNextWindow(ctx, task) {
\t\t\t\treturn ctx.Err() // context canceled during sleep
\t\t\t}
\t\t\t// reload task config after waking up
\t\t\tif updatedTask, err := GetTaskInfo(ctx, task.Id); err == nil && updatedTask != nil {
\t\t\t\ttask = updatedTask
\t\t\t}
\t\t}

\t\t// check pause status"""

if old_loop in content and "schedule check" not in content:
    content = content.replace(old_loop, new_loop, 1)
    changes += 1
    print("OK: Added schedule check in processTaskRecipients loop")
else:
    print("SKIP: processTaskRecipients schedule check")

# 5. Add send delay in processRecipientBatch (after rate control wait, before pool submit)
old_rate = """\t\t// wait for rate control
\t\tif err := e.rateController.Wait(ctx); err != nil {
\t\t\tif errors.Is(err, context.Canceled) {
\t\t\t\tsafeClose() // safe close channel
\t\t\t\treturn err
\t\t\t}
\t\t\t// record error but continue
\t\t\tg.Log().Debugf(ctx, "Rate limit wait error: %v", err)

\t\t}"""

new_rate = """\t\t// wait for rate control
\t\tif err := e.rateController.Wait(ctx); err != nil {
\t\t\tif errors.Is(err, context.Canceled) {
\t\t\t\tsafeClose() // safe close channel
\t\t\t\treturn err
\t\t\t}
\t\t\tg.Log().Debugf(ctx, "Rate limit wait error: %v", err)
\t\t}

\t\t// cold mail: delay between emails (anti-ban)
\t\tif task.SendDelay > 0 {
\t\t\tselect {
\t\t\tcase <-time.After(time.Duration(task.SendDelay) * time.Second):
\t\t\tcase <-ctx.Done():
\t\t\t\tsafeClose()
\t\t\t\treturn ctx.Err()
\t\t\t}
\t\t}"""

if old_rate in content and "cold mail: delay" not in content:
    content = content.replace(old_rate, new_rate, 1)
    changes += 1
    print("OK: Added send delay in processRecipientBatch")
else:
    print("SKIP: processRecipientBatch delay")

# 6. Patch sendEmail to use sender rotation
old_send = """\tcurrentTask := task
\tif e.taskConfig != nil {
\t\tcurrentTask = e.taskConfig
\t}

\t// get rendered content and subject
\trenderedContent, renderedSubject := e.personalizeEmail(ctx, content, currentTask, recipient)

\tsender, err := mail_service.NewEmailSenderWithLocal(currentTask.Addresser)"""

new_send = """\tcurrentTask := task
\tif e.taskConfig != nil {
\t\tcurrentTask = e.taskConfig
\t}

\t// sender rotation: get next sender from pool if configured
\tactiveSender := currentTask.Addresser
\tactiveName := currentTask.FullName
\tif e.senderRotation != nil {
\t\tpoolSender, err := e.senderRotation.GetNextSender(ctx)
\t\tif err != nil {
\t\t\tg.Log().Warningf(ctx, "Sender rotation failed: %v", err)
\t\t} else if poolSender != nil {
\t\t\tactiveSender = poolSender.Email
\t\t\tactiveName = poolSender.Name
\t\t}
\t}

\t// get rendered content and subject
\trenderedContent, renderedSubject := e.personalizeEmail(ctx, content, currentTask, recipient)

\tsender, err := mail_service.NewEmailSenderWithLocal(activeSender)"""

if old_send in content and "sender rotation: get next" not in content:
    content = content.replace(old_send, new_send, 1)
    changes += 1
    print("OK: Added sender rotation in sendEmail")
else:
    print("SKIP: sendEmail sender rotation")

# 7. Patch display name to use rotated sender name
old_name = """\t// set sender display name
\tif currentTask.FullName != "" {
\t\tmessage.SetRealName(currentTask.FullName)
\t}"""

new_name = """\t// set sender display name (use rotated sender name if available)
\tif activeName != "" {
\t\tmessage.SetRealName(activeName)
\t} else if currentTask.FullName != "" {
\t\tmessage.SetRealName(currentTask.FullName)
\t}"""

if old_name in content and "rotated sender name" not in content:
    content = content.replace(old_name, new_name, 1)
    changes += 1
    print("OK: Patched display name to use rotated sender")
else:
    print("SKIP: display name patch")

# 8. Track base URL by active sender (not task default)
old_baseurl = """\t//Tracking emails
\tbaseURL := domains.GetBaseURLBySender(currentTask.Addresser)"""

new_baseurl = """\t//Tracking emails
\tbaseURL := domains.GetBaseURLBySender(activeSender)"""

if old_baseurl in content and content.count(new_baseurl) == 0:
    content = content.replace(old_baseurl, new_baseurl, 1)
    changes += 1
    print("OK: BaseURL uses activeSender for tracking")
else:
    print("SKIP: baseURL tracking patch")

# 9. Increment daily counter after successful send
old_success = """\treturn &SendResult{
\t\tRecipientID: recipient.Id,
\t\tMessageID:   messageID,
\t\tSuccess:     true,
\t\tError:       nil,
\t}
}
"""

new_success = """\t// track daily send count for sender rotation
\tif e.senderRotation != nil {
\t\tIncrementSenderDailyCount(ctx, activeSender)
\t}

\treturn &SendResult{
\t\tRecipientID: recipient.Id,
\t\tMessageID:   messageID,
\t\tSuccess:     true,
\t\tError:       nil,
\t}
}
"""

if old_success in content and "track daily send count" not in content:
    content = content.replace(old_success, new_success, 1)
    changes += 1
    print("OK: Added daily counter increment after send")
else:
    print("SKIP: daily counter increment")

with open(path, "w", encoding="utf-8", newline="\n") as f:
    f.write(content)

print(f"\nTotal changes: {changes}")
