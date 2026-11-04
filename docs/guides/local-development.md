# Local development

Phase 1 runs the Next.js workspace directly and the Go API/PostgreSQL services through Docker Compose. Use only the fictional seed records provided for the prototype. This setup has no authentication and is not suitable for real patient information.

## Prerequisites

- Node.js supported by the pinned Next.js package and pnpm 12.10.1.
- Go 1.26 or newer in the 1.26 line for local API development.
- Docker Engine with Compose v2, able to pull the PostgreSQL image and build the API image.
- OpenSSL for generating local credentials.

## Prepare the local environment

From the repository root:

```bash
pnpm install --frozen-lockfile
docker compose -f deploy/compose.yaml up -d --build db api
docker compose -f deploy/compose.yaml --profile dev-seed run --rm seed
```

If `.env` or `secrets/` does not exist, run `./scripts/dev-secrets.sh` once before starting Compose. It creates non-secret local settings plus three random database credentials with restrictive permissions. It refuses to overwrite an existing secrets directory. Keep these files local and never paste their contents into logs, issues, or commits. Existing database roles keep the credentials with which the database volume was initialized; only reset a disposable local database using a planned reset procedure.

Compose starts PostgreSQL, applies migrations, then starts the Go API on `127.0.0.1:8080`. Seeding is explicit and repeatable; it adds only the fictional demo records. Check the services and API health:

```bash
docker compose -f deploy/compose.yaml ps
curl -fsS http://127.0.0.1:8080/health/live
curl -fsS http://127.0.0.1:8080/health/ready
```

Inspect the API-backed dashboard counts and episode rows:

```bash
curl -fsS http://127.0.0.1:8080/api/v1/dashboard/overview | python3 -m json.tool
```

Copy an `id` from the returned `episodes` array and request its detail:

```bash
curl -fsS http://127.0.0.1:8080/api/v1/episodes/EPISODE_ID | python3 -m json.tool
```

In a second terminal, run the web app and open <http://localhost:3000>:

```bash
pnpm dev
```

The Next.js server reads `RELAY_API_URL` from `.env`; the default is `http://127.0.0.1:8080`. Select a patient name in the dashboard queue to open its episode details. The UI renders API data only. If the API is unavailable, it shows an unavailable state instead of fallback counts or rows.

## Smoke checks

- `GET /health/live` returns `{"status":"ok"}` when the API process is running.
- `GET /health/ready` returns `{"status":"ready"}` when PostgreSQL is reachable.
- `GET /api/v1/dashboard/overview` returns seeded fictional summary counts and open episode rows.
- `GET /api/v1/episodes/<id>` returns episode detail. An unknown or closed episode returns HTTP 404.
- In the browser, episode names open details and the light/dark theme control works on both routes.
- To see the unavailable state, run `docker compose -f deploy/compose.yaml stop api`, reload the page, then restore the API with `docker compose -f deploy/compose.yaml start api`.
- To see the not-found UI, navigate to `/episodes/not-an-id` in the browser.

These routes are unauthenticated. Use only the included fictional records; do not enter real patient information or expose the services beyond loopback.

## Available commands

| Command | Purpose |
| --- | --- |
| `pnpm dev` | Run the Next.js workspace. |
| `pnpm build` | Build the web application. |
| `pnpm start` | Run the built web application. |
| `pnpm lint` | Run ESLint. |
| `go run ./backend/cmd/api migrate` | Apply embedded SQL migrations (requires database environment variables and migrator credentials). |
| `go run ./backend/cmd/api seed` | Add fictional demo records; allowed only when `APP_ENV=development`. |
| `docker compose -f deploy/compose.yaml logs -f api` | Follow local API logs. |

The API image waits for database migration before startup. The seed command is an explicit, repeatable development step and is not run automatically by `up`.

## Next.js documentation rule

Before changing application code, follow the repository `AGENTS.md`: read the relevant guide under `node_modules/next/dist/docs/` for the installed version.

## Current limits

The API is read-only and unauthenticated. Do not expose it beyond loopback or enter real patient information. Production use needs identity and access controls, TLS and managed secrets, institutional privacy/security approval, data retention and backup policies, clinical workflow review, and a deployment threat model. See the [security and secrets guide](security-secrets.md).
