#!/bin/sh
set -eu

admin_password=$(cat /run/secrets/postgres_admin_password)
APP_DB_PASSWORD=$(cat /run/secrets/app_db_password)
MIGRATOR_DB_PASSWORD=$(cat /run/secrets/migrator_db_password)
export APP_DB_PASSWORD MIGRATOR_DB_PASSWORD

passfile=$(mktemp)
chmod 600 "$passfile"
trap 'rm -f "$passfile"' EXIT HUP INT TERM
printf '%s:%s:%s:%s:%s\n' "$PGHOST" "$PGPORT" "$POSTGRES_DB" "$POSTGRES_USER" "$admin_password" > "$passfile"
unset admin_password
export PGPASSFILE="$passfile"

psql \
	--host "$PGHOST" \
	--username "$POSTGRES_USER" \
	--dbname "$POSTGRES_DB" \
	--set ON_ERROR_STOP=1 <<'SQL'
\getenv app_password APP_DB_PASSWORD
\getenv migrator_password MIGRATOR_DB_PASSWORD

SELECT format('CREATE ROLE renalrelay_app LOGIN PASSWORD %L', :'app_password')
WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'renalrelay_app')
\gexec
SELECT format('CREATE ROLE renalrelay_migrator LOGIN PASSWORD %L', :'migrator_password')
WHERE NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'renalrelay_migrator')
\gexec
GRANT CONNECT ON DATABASE renalrelay TO renalrelay_app, renalrelay_migrator;
GRANT USAGE ON SCHEMA public TO renalrelay_app, renalrelay_migrator;
GRANT CREATE ON SCHEMA public TO renalrelay_migrator;
ALTER DEFAULT PRIVILEGES FOR ROLE renalrelay_migrator IN SCHEMA public
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO renalrelay_app;
ALTER DEFAULT PRIVILEGES FOR ROLE renalrelay_migrator IN SCHEMA public
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO renalrelay_app;
SQL

unset APP_DB_PASSWORD MIGRATOR_DB_PASSWORD PGPASSFILE
