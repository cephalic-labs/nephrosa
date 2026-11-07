# Single-instance demo deployment

This deployment runs Next.js, the Go API, and PostgreSQL with Docker Compose. Only the web port is published. The separate `renal-relay-demo` Compose project has its own database volume; it does not use the local development stack's volume. Use fictional records only.

## Prepare

Install Docker Engine with Compose 2.24.4 or newer, Make, and OpenSSL on the instance. Clone the repository and, from its root, generate the ignored `.env` and owner-only secret files:

```sh
./scripts/dev-secrets.sh
```

The script refuses to overwrite an existing `secrets/` directory. Preserve the secret files along with the database volume; do not regenerate credentials for an existing database. Set `LOCAL_UID` and `LOCAL_GID` in `.env` to the UID and GID of the user who owns the secret files if they are not already correct. Set `WEB_PORT` to the desired host port (default `3000`) and keep `APP_ENV=development`: seeding requires it and the database connection inside the Compose network is not encrypted.

For calls, configure the Sarvam settings in `.env`, put the workspace API key in `secrets/sarvam_api_key`, and set `SARVAM_WEBHOOK_URL` to the public HTTPS base ending in `/api/v1/sarvam/webhooks`. Keep `SARVAM_WEBHOOK_ENABLED=false` initially. See the [Samvaad guide](../docs/guides/samvaad-demo.md) for agent setup and credentials. The generated demo account passwords are in `secrets/demo_clinician_password` and `secrets/demo_coordinator_password`; the usernames are `clinician` and `coordinator`.

## Start and operate

```sh
make up       # build and start the stack, roles, and migrations
make seed     # explicitly provision fictional records and demo accounts
make status   # inspect containers, including one-shot setup jobs
make logs     # follow logs; Ctrl-C stops following
make down     # stop containers, retaining the database volume
```

`make build` builds the images without starting them. Run `make up` after code or configuration changes; run `make seed` only when needed, since reseeding rotates demo credentials and invalidates existing sessions. For updates, pull the desired code and run `make up` again. Do not run `docker compose down --volumes` unless intentionally deleting all demo data. Back up or preserve the volume and associated secret files if the demo data must be retained.

## HTTPS and callbacks

Point your HTTPS domain or tunnel at the published web port. Sign-in requires HTTPS because the production frontend sets secure session cookies. Configure the HTTPS proxy or tunnel to omit or redact the complete callback path and body from access logs and analytics; callback URLs contain secrets. Verify this and HTTPS routing before setting `SARVAM_WEBHOOK_ENABLED=true` and running `make up` to recreate the API. Check logs from the proxy, tunnel, and application without publishing callback URLs or bodies.

Next.js forwards only `/api/v1/sarvam/webhooks/:secret` to the internal API. Verify a request with an invalid callback secret is rejected and other API routes are not publicly reachable through the frontend. Then use a consenting demo phone for a live call and confirm the callback updates the episode. Consult the [callback security and recovery instructions](../docs/guides/samvaad-demo.md#callback-security-and-recovery) for operational details.
