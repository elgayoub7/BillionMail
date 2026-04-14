# Security Policy

## Reporting a Vulnerability

Please report security issues to `1249648969@qq.com`

---

# Security Audit Report - Hardened Fork

This document describes all security vulnerabilities found in the original BillionMail repository and the fixes applied in this hardened fork.

**Audit Date:** 2026-04-14
**Auditor:** Manual code review
**Severity Scale:** Critical (Red) / Medium (Yellow) / Low (Green)

---

## Critical Vulnerabilities Fixed (8)

### 1. JWT Secret Derived from Database Passwords
- **File:** `core/internal/service/rbac/jwt.go:36-56`
- **Issue:** JWT signing secret was constructed by concatenating `DBPASS + REDISPASS`. If the database is compromised, JWT tokens can be forged.
- **Fix:** Added dedicated `JWT_SECRET` environment variable. Falls back to old behavior with a warning log if not set.
- **Action Required:** Set `JWT_SECRET` in your `.env` file:
  ```bash
  JWT_SECRET=$(openssl rand -hex 32)
  ```

### 2. API Tokens Without Expiration
- **File:** `core/internal/service/rbac/jwt.go:110-133`
- **Issue:** `GenerateApiToken` created JWTs with `ExpiresAt: nil` — tokens never expired. A leaked token would be compromised indefinitely.
- **Fix:** API tokens now expire after 365 days (configurable via `jwt.apiExpiry` in config).
- **Impact:** Existing API tokens will continue to work (backward compatible). New tokens will have an expiry.

### 3. Default Credentials Committed in Repository
- **File:** `env_init`
- **Issue:** The file contained actual passwords (`DBPASS=NauF7ysRYyt9HTOiOn4JjIAL3QcRZnzj`, `ADMIN_PASSWORD=billion`).
- **Fix:** All credentials replaced with `CHANGE_ME` placeholders. The `install.sh` script already generates random passwords for new installs — this only affected manual Docker deployments.

### 4. Container Running as Root
- **File:** `Dockerfiles/core/Dockerfile`, `Dockerfiles/core/supervisord.conf`
- **Issue:** No `USER` directive in Dockerfile. All supervisord processes ran as `root`.
- **Fix:** Created non-root user `billionmail`. The core application now runs as `billionmail`. Fail2ban and crond still run as root (required for their functionality).

### 5. Docker Socket Mounted with Elevated Capabilities
- **File:** `docker-compose.yml:195-203`
- **Issue:** `/var/run/docker.sock` was mounted read-only into the container, which also had `NET_ADMIN` and `NET_RAW` capabilities. An exploit in the application could gain near-root access to the host.
- **Fix:** Docker socket mount removed. `NET_ADMIN` and `NET_RAW` capabilities removed. Only `NET_BIND_SERVICE` is retained (needed for binding ports < 1024).

### 6. Mail Ports Exposed Without IP Restriction
- **File:** `docker-compose.yml:134-136, 91-94`
- **Issue:** SMTP (25), SMTPS (465), Submission (587), IMAP (143), IMAPS (993), POP3 (110), POP3S (995) were all bound to `0.0.0.0`, creating a potential open relay.
- **Fix:** All mail ports now default to `127.0.0.1` binding. Only HTTP/HTTPS remain publicly accessible. To expose mail ports publicly, override in `.env`:
  ```bash
  SMTP_PORT=25          # public (0.0.0.0)
  SMTP_PORT=127.0.0.1:25  # localhost only
  ```

### 7. TLS Verification Globally Disabled
- **File:** `core/internal/service/public/common.go:276-286`
- **Issue:** `InsecureSkipVerify: true` was set on ALL outgoing HTTP connections, enabling Man-in-the-Middle attacks on webhooks, API calls, and tracking requests.
- **Fix:** `GetHttpClient()` now uses proper TLS verification with TLS 1.2+ minimum. A separate `GetInsecureHttpClient()` is available for internal connections where self-signed certificates are used.

### 8. SSRF via AI Chat
- **File:** `core/internal/service/askai/openai.go:176-204`
- **Issue:** The AI chat's HTTP request tool accepted arbitrary URLs from LLM responses (potential prompt injection), with no validation. Combined with `InsecureSkipVerify`, this allowed Server-Side Request Forgery attacks against internal services.
- **Fix:** Added `isUrlSafe()` validation function that blocks:
  - Private/internal IP ranges (RFC 1918)
  - Loopback addresses (`localhost`, `127.0.0.1`, `::1`)
  - Link-local addresses
  - Cloud metadata endpoints (`169.254.169.254`, `metadata.google.internal`)
  - Non-HTTP/HTTPS schemes

---

## Medium Vulnerabilities Fixed (1)

### 9. Weak DKIM Key Size
- **File:** `install.sh:1044`
- **Issue:** DKIM keys were generated with 1024-bit strength, below the 2048-bit minimum recommended by RFC 8301.
- **Fix:** Changed to `rspamadm dkim_keygen -b 2048`.
- **Note:** Existing keys are NOT affected. New installations and new domains will use 2048-bit keys.

---

## Remaining Items (to be addressed in future updates)

| # | Issue | Severity | Status |
|---|---|---|---|
| 1 | Refresh token does not invalidate the old one | Medium | Planned |
| 2 | API token stored in plaintext in DB | Medium | Planned |
| 3 | Password complexity checks commented out | Medium | Planned |
| 4 | `math/rand` instead of `crypto/rand` for `RandomStr()` | Medium | Planned |
| 5 | CORS middleware using default (permissive) settings | Medium | Planned |
| 6 | gmaps-scraper exposed on port 8080 without auth | Medium | Planned |
| 7 | No rate limiting on public API endpoints | Medium | Planned |
| 8 | Session cookie missing Secure/HttpOnly/SameSite flags | Medium | Planned |
| 9 | Sender address spoofing possible via API | Medium | Planned |
| 10 | XSS risk in email template HTML content | Medium | Planned |

---

## Migration Guide

See [MIGRATION.md](MIGRATION.md) for step-by-step instructions on how to apply these security fixes to an existing BillionMail installation.
