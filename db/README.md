# Database scaffold

PostgreSQL is the local prototype store. The migrations establish workflow entities, integrity constraints, hashed demo sessions, patient contact numbers, and correlated call attempts. The seed provides two fictional monitoring episodes and demo-account provisioning. The schema supports the local fictional-data demonstration.

## Migrations

- Add a new, ordered `NNNNNN_description.sql` file under `migrations/` for schema changes. Applied migrations are immutable; use a follow-up migration to correct them.
- `migrations/embed.go` embeds SQL into the API binary.
- The `migrate` command executes each unapplied migration transactionally under the separate `renalrelay_migrator` role and records it in `schema_migrations`.
- Review foreign keys, constraints, indexes, permissions, retention, and audit implications with every schema change.

## Demo seed

- `seed/001_fictional_demo_data.sql` contains generated fictional names/references and illustrative workflow entries only.
- The seed uses stable IDs and conflict-safe inserts for repeatable local use.
- `seed/embed.go` embeds SQL into the API binary. Run it explicitly with `APP_ENV=development` after migrations; Compose does not run it automatically.
- The seed command reads the clinician/coordinator password files, hashes the passwords, and revokes older demo sessions.
- Never add real patient identifiers, laboratory results, contact details, or imported data to the seed.

## Roles

The `roles` Compose bootstrap service creates the migration role with schema creation rights and the application role with limited table/sequence DML. It runs as the host UID/GID that owns the owner-only local secret files, then connects to PostgreSQL with the administrator credential. PostgreSQL itself mounts only its administrator secret. The API uses the application role for reads and transactional workflow writes. Clinician/coordinator authorization is enforced in the Go service. Credentials are sourced from ignored `secrets/` files as described in the [security guide](../docs/guides/security-secrets.md).
