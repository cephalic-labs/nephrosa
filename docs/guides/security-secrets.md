# Local secrets and deployment credentials

## Local development

Local credentials are generated into the ignored `secrets/` directory by `scripts/dev-secrets.sh`. The script creates separate random PostgreSQL administrator, application, and migration passwords; it never prints their values. Secret files are created with owner-only permissions, and the generated `.env` contains non-secret development settings only.

The Compose stack mounts these files as secrets. The API connects as the restricted `renalrelay_app` role. The one-shot migration container uses `renalrelay_migrator`. The PostgreSQL administrator credential is used only by the database initialization process.

Never commit `secrets/`, `.env`, credential files, or screenshots/logs containing secret values. Do not copy production credentials into local development. The local Compose setup is a development boundary, not production secret storage.

## Production

Use the deployment platform's managed secret store and inject credentials at runtime. Use separate least-privilege database credentials for the application and migrations. Require TLS to managed databases, rotate credentials through the secret provider, and ensure backups and diagnostic logs exclude secret values.

Never put server credentials in `NEXT_PUBLIC_*` variables; Next.js inlines those values into browser JavaScript. Keep browser configuration limited to values explicitly intended to be public.

## Rotation

Regenerate local credentials by stopping the stack and removing only disposable development secrets and database volumes, then rerun `scripts/dev-secrets.sh`. Do not rotate or delete credentials for an environment containing data that must be retained without a planned database credential change and recovery procedure.
