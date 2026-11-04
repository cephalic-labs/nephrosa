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
./scripts/dev-secrets.sh
docker compose -f deploy/compose.yaml up -d --build db api
docker compose -f deploy/compose.yaml --profile dev-seed run --rm seed
```

The script creates `.env` and three random database credentials under the ignored `secrets/` directory. It will not overwrite an existing secrets directory. Keep those files local and never paste their contents into logs, issues, or commits. If resetting the disposable local database, stop the stack, remove the `postgres_data` volume, archive or remove `secrets/`, then run the script again; existing database roles retain the credentials with which the volume was initialized.

In a second terminal:

```bash
pnpm dev
```

Open <http://localhost:3000>. The page reads the dashboard API at `RELAY_API_URL` (default `http://127.0.0.1:8080`). It displays clearly labeled fictional sample data if the API is unavailable. To inspect service state:

```bash
docker compose -f deploy/compose.yaml ps
curl http://127.0.0.1:8080/health/live
curl http://127.0.0.1:8080/health/ready
```

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
