# Migration Guide - Security Hardening

This guide explains how to migrate from the original BillionMail to this security-hardened fork.

---

## For New Installations

No special steps needed. The hardened fork works identically to the original:

```bash
cd /opt && git clone https://github.com/YOUR_FORK/BillionMail && cd BillionMail && bash install.sh
```

The `install.sh` script already generates random passwords. The new `JWT_SECRET` will be auto-generated.

---

## For Existing Installations

### Step 1: Backup

```bash
# Backup your data
cp -r /opt/BillionMail /opt/BillionMail-backup-$(date +%Y%m%d)
cp /opt/BillionMail/.env /opt/BillionMail-backup-$(date +%Y%m%d)/.env.backup
```

### Step 2: Pull Changes

```bash
cd /opt/BillionMail
git stash                    # Save any local changes
git remote add fork https://github.com/YOUR_FORK/BillionMail
git fetch fork
git checkout fork/main       # Or whichever branch
git stash pop                # Restore local changes if any
```

### Step 3: Update `.env` File

Add the new `JWT_SECRET` variable:

```bash
# Generate a secure JWT secret
echo "JWT_SECRET=$(openssl rand -hex 32)" >> /opt/BillionMail/.env
```

### Step 4: Rebuild Docker Images

```bash
cd /opt/BillionMail
docker compose build core-billionmail   # Rebuild with new Dockerfile
docker compose up -d
```

### Step 5: Verify

```bash
# Check all containers are running
docker compose ps

# Verify the core is NOT running as root
docker exec core-billionmail-1 whoami
# Expected output: billionmail

# Verify mail ports are on localhost only
docker compose port postfix-billionmail 25
# Expected output: 127.0.0.1:25

# Verify JWT_SECRET is set
docker exec core-billionmail-1 grep JWT_SECRET /opt/billionmail/.env
```

### Step 6: (Optional) Rotate API Tokens

Since API tokens now have an expiration, existing tokens continue to work. To generate a new token with expiry:

1. Log into BillionMail admin panel
2. Go to API settings
3. Generate a new API token
4. Update your integrations with the new token

---

## Configuration Changes

### New Environment Variable

| Variable | Default | Description |
|---|---|---|
| `JWT_SECRET` | (falls back to DBPASS+REDISPASS) | Dedicated secret for JWT signing. Generate with `openssl rand -hex 32` |
| `jwt.apiExpiry` | `31536000` (365 days in seconds) | API token expiration in seconds. Set in YAML config |

### Port Binding Changes

Mail ports now default to localhost-only binding. To restore public access (if you have a proper firewall):

```bash
# In your .env file, change:
SMTP_PORT=25              # Public (original behavior)
SMTP_PORT=127.0.0.1:25    # Localhost only (new default)

# Available options:
# 25                → 0.0.0.0:25 (public)
# 127.0.0.1:25      → localhost only
# 192.168.1.100:25  → specific IP only
```

### Removed Docker Socket

If you were relying on the Docker socket for container management from within BillionMail:

```bash
# To re-enable (NOT recommended for production):
# In docker-compose.yml, uncomment:
# - /var/run/docker.sock:/var/run/docker.sock:ro
```

---

## Breaking Changes

| Change | Impact | Mitigation |
|---|---|---|
| Mail ports on 127.0.0.1 | External mail clients can't connect directly | Use a reverse proxy (nginx) to forward mail ports, or override in `.env` |
| Docker socket removed | Container management features from admin panel won't work | Re-enable socket only if needed and you understand the risk |
| `NET_ADMIN`/`NET_RAW` removed | Firewall management from admin panel limited | Run firewall commands directly on the host |

---

## Troubleshooting

### "permission denied" errors after upgrade

The core process now runs as `billionmail` user. Fix permissions:

```bash
docker exec core-billionmail-1 chown -R billionmail:billionmail /opt/billionmail/core/data
docker exec core-billionmail-1 chown -R billionmail:billionmail /opt/billionmail/logs
```

### Existing API token stopped working

Tokens created before this update have no expiry and continue to work. If you see auth errors:

1. Check that `JWT_SECRET` is properly set in `.env`
2. Generate a new API token from the admin panel
3. Verify the token is being sent correctly in the `Authorization` header

### TLS errors on internal connections

If you have self-signed certificates for internal services and see TLS errors:

This is the intended behavior of the security fix. If you need to connect to a service with a self-signed certificate, use `GetInsecureHttpClient()` in Go code, but only for internal/trusted connections.
