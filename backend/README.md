# Workflow API scaffold

The Go service currently provides a minimal read path for the prototype dashboard. It has no authentication, authorization, write endpoints, or production readiness claim.

## Commands

The `backend/cmd/api` executable accepts:

| Command | Behavior |
| --- | --- |
| `serve` (default) | Start the HTTP API after connecting to PostgreSQL. |
| `migrate` | Apply ordered embedded SQL migrations using the migration database role. |
| `seed` | Apply fictional seed records; refuses to run unless `APP_ENV=development`. |
| `healthcheck` | Probe the local liveness endpoint for container health. |

The container builds as a static binary and runs as a non-root user. Compose supplies the password-file path at runtime; the password itself is not passed as an environment value.

## HTTP surface

| Method and path | Purpose |
| --- | --- |
| `GET /health/live` | Process liveness; does not require a database query. |
| `GET /health/ready` | Database readiness. |
| `GET /api/v1/dashboard/overview` | Counts and up to ten open episode summaries. |
| `GET /api/v1/episodes/:id` | Read-only details for an open episode: test orders and verification, contacts, barriers, tasks, report receipt metadata, clinician acknowledgement and next action, and audit event labels. |

Episode details are read in a read-only repeatable-read transaction. Report contents and audit metadata are not returned. A missing or closed episode returns 404; a database failure returns a generic 503 response. Responses include `Cache-Control: no-store` and basic browser security headers. Errors are returned with a generic message while internal details are logged server-side. Do not log credentials, conversation contents, or unnecessary patient details.

The route has no authentication or authorization. It is for the local fictional-data prototype only and must not be exposed to real-patient data or an untrusted network.

## Configuration

`APP_ENV`, `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_NAME`, `DATABASE_USER`, and `DATABASE_PASSWORD_FILE` configure startup. `DATABASE_SSLMODE` defaults to `disable` only for development; other environments must set it, and production requires `verify-full`. `HTTP_ADDR` defaults to `:8080`.

Before changing this service, review [the product and workflow boundaries](../PRODUCT.md), [architecture baseline](../docs/architecture/renal-relay-hld-lld.md), and [secret handling](../docs/guides/security-secrets.md). Validate configuration, database role permissions, and errors without printing secret values.
