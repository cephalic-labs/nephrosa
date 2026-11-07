# Workflow API

The Go service owns authenticated demo sessions, workflow mutations and audit writes, dashboard/detail reads, and the Samvaad instant-calling adapter. See [the demo guide](../docs/guides/samvaad-demo.md). The running API, real Samvaad call and callback, and fictional workflow closeout were verified on 7 November 2026.

## Commands

The `backend/cmd/api` executable accepts:

| Command | Behavior |
| --- | --- |
| `serve` (default) | Start the HTTP API after connecting to PostgreSQL. |
| `migrate` | Apply ordered embedded SQL migrations using the migration database role. |
| `seed` | Apply fictional seed records; refuses to run unless `APP_ENV=development`. |
| `reset-demo --confirm-fictional-reset` | Truncate fictional workflow records using the migration role after development/database/reference checks. |
| `healthcheck` | Probe the local liveness endpoint for container health. |

The container builds as a static binary and runs as a non-root user. Compose supplies the password-file path at runtime; the password itself is not passed as an environment value.

## HTTP surface

| Method and path | Purpose |
| --- | --- |
| `GET /health/live` | Process liveness; does not require a database query. |
| `GET /health/ready` | Database readiness. |
| `GET /api/v1/dashboard/overview` | Counts and up to 100 episode summaries. |
| `GET /api/v1/episodes/:id` | Details for an open or closed episode: test orders and verification, contacts, barriers, tasks, report receipt metadata, clinician acknowledgement and next action, and audit event labels. |

Episode details are read in a read-only repeatable-read transaction. Report references and human resolution notes are returned. Report document contents and other audit metadata stay outside the response. A missing episode returns 404; a database failure returns a generic 503 response. Responses include `Cache-Control: no-store` and basic browser security headers. Errors are returned with a generic message while internal details are logged server-side. Do not log credentials, conversation contents, or unnecessary patient details.

All care-team endpoints require a bearer session token. The web server holds that token in an HTTP-only cookie and forwards it server-side. Only health, sign-in, and the secret callback URL are unauthenticated by session. Expose only the callback prefix publicly; use fictional records.

## Configuration

`APP_ENV`, `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_NAME`, `DATABASE_USER`, and `DATABASE_PASSWORD_FILE` configure startup. `DATABASE_SSLMODE` defaults to `disable` only for development; other environments must set it, and production requires `verify-full`. `HTTP_ADDR` defaults to `:8080`.

Before changing this service, review [the product and workflow boundaries](../PRODUCT.md), [architecture baseline](../docs/architecture/renal-relay-hld-lld.md), and [secret handling](../docs/guides/security-secrets.md). Validate configuration, database role permissions, and errors without printing secret values.

## Workflow API

- `POST /api/v1/sessions`: login/password → opaque eight-hour token and user.
- `GET /api/v1/session`, `DELETE /api/v1/session`: current user and revocation.
- `GET /api/v1/users`: active care-owner choices.
- `POST /api/v1/episodes`: clinician creates patient and episode atomically, including ordered tests and consent provenance.
- `POST /api/v1/episodes/:id/actions`: consent, assign (clinician), verify_test, report, resolve_task, resolve_barriers, review (clinician), close (clinician).
- `POST /api/v1/episodes/:id/calls`, `GET /api/v1/episodes/:id/calls`: start and inspect real call attempts.
- `POST /api/v1/episodes/:id/calls/:call/reconcile`: record provider-log verification before releasing an outstanding call.
- `POST /api/v1/sarvam/webhooks/:secret`: bounded instant-call callback authenticated by a per-call capability. Disabled by default.

Use `?history=true` on dashboard overview to include closed episodes. The list is capped at 100; summary counts cover the full dataset. Episode details include current consent and contact information; audit entries include human resolution notes. The initial append-only trigger protects audit updates/deletes; only the explicit migrator-only development reset truncates fictional history.
