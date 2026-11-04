# Local service orchestration

`compose.yaml` runs PostgreSQL and the Go API on loopback for local development. It is not a production deployment manifest. The web application runs separately through `pnpm dev`.

## Start

From the repository root, install dependencies and generate local secrets once:

```bash
pnpm install --frozen-lockfile
./scripts/dev-secrets.sh
docker compose -f deploy/compose.yaml up -d --build db api
docker compose -f deploy/compose.yaml --profile dev-seed run --rm seed
pnpm dev
```

The API waits for PostgreSQL readiness and successful migrations. The fictional seed is an explicit development-profile command. The default ports are 3000 for Next.js, 8080 for the API, and 5432 for PostgreSQL; API and database ports bind to `127.0.0.1`.

## Stop and reset

Stop without deleting the local database:

```bash
docker compose -f deploy/compose.yaml down
```

The `postgres_data` volume persists across restarts. For a disposable reset, stop the stack and run `docker compose -f deploy/compose.yaml down --volumes`; this permanently removes local database data. Keep the existing secret files unless intentionally resetting credentials together with the volume.

## Security properties and limits

- Containers drop Linux capabilities, use read-only root filesystems where practical, and disable privilege escalation.
- PostgreSQL uses SCRAM password authentication. App and migration credentials are separate Docker secrets backed by local ignored files.
- Ports bind to loopback; do not change this for a shared or public environment.
- There is no API authentication or workflow write surface. Never use real patient information.
- Review [secret handling](../docs/guides/security-secrets.md) before changing credentials or deployment configuration.
