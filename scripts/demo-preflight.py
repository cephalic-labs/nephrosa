#!/usr/bin/env python3
"""Check local demo configuration without printing credentials or config values."""
import os
import re
from pathlib import Path
from urllib.parse import urlparse

ROOT = Path(__file__).resolve().parent.parent
REQUIRED = ("SARVAM_ORG_ID", "SARVAM_WORKSPACE_ID", "SARVAM_APP_ID", "SARVAM_APP_VERSION", "SARVAM_CONNECTION_ID", "SARVAM_AGENT_PHONE_NUMBER", "SARVAM_WEBHOOK_URL", "DEMO_CALL_PHONE_NUMBER")

def check():
    settings = {}
    env = ROOT / ".env"
    if env.exists():
        for line in env.read_text().splitlines():
            if line.strip() and not line.lstrip().startswith("#") and "=" in line:
                key, value = line.split("=", 1)
                settings[key.strip()] = value.strip().strip("\"'")
    settings.update(os.environ)
    problems = [f"Missing {key} in root .env" for key in REQUIRED if not settings.get(key)]
    for key in ("SARVAM_AGENT_PHONE_NUMBER", "DEMO_CALL_PHONE_NUMBER"):
        if settings.get(key) and not re.fullmatch(r"\+[1-9][0-9]{7,14}", settings[key]):
            problems.append(f"Invalid E.164 format: {key}")
    if not settings.get("SARVAM_APP_VERSION", "").isdigit() or int(settings.get("SARVAM_APP_VERSION") or "0") < 1:
        problems.append("SARVAM_APP_VERSION must be a positive committed version")
    url = urlparse(settings.get("SARVAM_WEBHOOK_URL", ""))
    if url.scheme != "https" or not url.netloc or url.query or url.fragment or not url.path.endswith("/api/v1/sarvam/webhooks"):
        problems.append("SARVAM_WEBHOOK_URL must be HTTPS and end in /api/v1/sarvam/webhooks")
    if settings.get("SARVAM_WEBHOOK_ENABLED") != "true":
        problems.append("SARVAM_WEBHOOK_ENABLED must be true after reviewing callback URL security")
    for name in ("sarvam_api_key", "demo_clinician_password", "demo_coordinator_password", "app_db_password", "migrator_db_password", "postgres_admin_password"):
        path = ROOT / "secrets" / name
        if not path.is_file() or not path.read_text().strip():
            problems.append(f"Missing or empty secrets/{name}")
        elif path.stat().st_mode & 0o077:
            problems.append(f"secrets/{name} must have owner-only permissions (chmod 600)")
    return problems

if __name__ == "__main__":
    problems = check()
    for problem in problems:
        print(problem)
    print("Configuration requires attention." if problems else "Local configuration is ready; no call has been placed. Check database/API health separately.")
    raise SystemExit(bool(problems))
