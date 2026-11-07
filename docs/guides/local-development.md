# Local development

Run Next.js on the host and the Go API with PostgreSQL through Docker Compose. The application provides authenticated clinician/coordinator demonstration accounts, audited workflow actions, and optional Sarvam Samvaad calls. Use fictional clinical records and a consenting demonstration phone.

## Prerequisites

- Node.js supported by the pinned Next.js package and pnpm 12.10.1.
- Go 1.26 or newer for direct API development and backend tests.
- Docker Engine with Compose, with access to the public PostgreSQL and Go images.
- OpenSSL for generating local credentials and Python 3 for the configuration preflight.

## Prepare and start

From the repository root, install dependencies and generate credentials on a fresh checkout.

```sh
pnpm install --frozen-lockfile
./scripts/dev-secrets.sh
```

The setup script creates owner-only database and demo-account passwords, an empty Samvaad key file, and the root `.env`. It refuses to overwrite an existing `secrets/` directory. On an existing installation, create only missing files and preserve credentials associated with the database volume. The [Samvaad guide](samvaad-demo.md) describes each required setting and the value-entry checkpoint before live testing.

Start the database, role bootstrap, migrations, API, and fictional seed.

```sh
docker compose --env-file .env -f deploy/compose.yaml up -d --build db api
docker compose --env-file .env -f deploy/compose.yaml --profile dev-seed run --rm seed
pnpm dev
```

Open <http://localhost:3000/login>. Use the `clinician` or `coordinator` login and the matching password in `secrets/demo_clinician_password` or `secrets/demo_coordinator_password`. Seeding hashes these passwords and revokes prior demo-account sessions. The PostgreSQL volume persists across restarts.

The explicit `--env-file .env` supplies root settings to Compose. The web server reads `RELAY_API_URL` from its Next.js environment, with `.env.local` taking precedence over `.env`. Its default is `http://127.0.0.1:8080`. Compose passes the Go service's environment and mounts its credentials separately.

## Verify the services

```sh
docker compose --env-file .env -f deploy/compose.yaml ps
curl -fsS http://127.0.0.1:8080/health/live
curl -fsS http://127.0.0.1:8080/health/ready
go test ./...
pnpm lint
pnpm build
```

The backend integration test requires `TEST_DATABASE_URL` for a disposable PostgreSQL database and a role allowed to create schemas. It creates and drops an isolated schema. Without that setting, Go explicitly skips the database test while running unit tests.

The dashboard, episode, and workflow endpoints require a bearer session token. Next.js holds that token in an HTTP-only cookie and forwards it server-side. The browser guide in [the demo runbook](samvaad-demo.md) exercises the authenticated path without exposing tokens in terminal output.

A clinician can create an episode, assign an owner, acknowledge its report, and authorize closure. Both care-team roles can record consent, request calls, resolve tasks/barriers, verify test completion, and record report receipt. Closed episodes remain available through dashboard history. Unknown episodes return HTTP 404.

## Live calls and configuration changes

Configure the agent, connected phone number, API key, public HTTPS callback, and `DEMO_CALL_PHONE_NUMBER` using the [Samvaad runbook](samvaad-demo.md). Run `python3 scripts/demo-preflight.py` after the owner has populated the files and confirmed readiness. This check validates local configuration without placing a call.

After changing the API key, app version, or callback settings, recreate the API to load the new values.

```sh
docker compose --env-file .env -f deploy/compose.yaml up -d --force-recreate api
```

The demo number is used when creating a fictional patient. Changing it in `.env` takes effect in a newly created test episode. Existing patients retain their saved numbers. Calls originate through the Go API after consent and duplicate-call checks.

## Stop and reset

Stop Next.js with Ctrl-C and stop the Compose services with the following command.

```sh
docker compose --env-file .env -f deploy/compose.yaml down
```

For an explicit fictional-data reset, follow the guarded reset procedure in the Samvaad guide. Keep the volume and credentials together when retaining records.

## Current limits

Calls start in English and contact history currently records English. Multilingual conversation accuracy and a dedicated live structured-barrier result remain unverified. The dashboard list is capped at 100 records. Webhooks use per-call secret URLs, and deployment must protect those URLs from disclosure. Review [secret handling](security-secrets.md) before changing credentials or deployment configuration.

Before changing Next.js code, read the installed guides under `node_modules/next/dist/docs/` as required by `AGENTS.md`.
