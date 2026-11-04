# Database scaffold

PostgreSQL is the local prototype store. The current migration establishes workflow entities and integrity constraints, while the seed provides two fictional monitoring episodes. The schema is a starting point for review, not a production clinical record model.

## Migrations

- Add a new, ordered `NNNNNN_description.sql` file under `migrations/` for schema changes. Applied migrations are immutable; use a follow-up migration to correct them.
- `migrations/embed.go` embeds SQL into the API binary.
- The `migrate` command executes each unapplied migration transactionally under the separate `renalrelay_migrator` role and records it in `schema_migrations`.
- Review foreign keys, constraints, indexes, permissions, retention, and audit implications with every schema change.

## Demo seed

- `seed/001_fictional_demo_data.sql` contains generated fictional names/references and illustrative workflow entries only.
- The seed uses stable IDs and conflict-safe inserts for repeatable local use.
- `seed/embed.go` embeds SQL into the API binary. Run it explicitly with `APP_ENV=development` after migrations; Compose does not run it automatically.
- Never add real patient identifiers, laboratory results, contact details, or imported data to the seed.

## Roles

The PostgreSQL initialization script creates an administrator for local initialization, a migration role with schema creation rights, and an application role with limited table/sequence DML. The current API only reads data; the application role is not intended as a clinical authorization boundary. Credentials are sourced from ignored `secrets/` files as described in the [security guide](../docs/guides/security-secrets.md).
