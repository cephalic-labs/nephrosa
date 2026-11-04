#!/bin/sh
set -eu

APP_DB_PASSWORD=$(cat /run/secrets/app_db_password)
MIGRATOR_DB_PASSWORD=$(cat /run/secrets/migrator_db_password)
export APP_DB_PASSWORD MIGRATOR_DB_PASSWORD

psql \
  --username "$POSTGRES_USER" \
  --dbname "$POSTGRES_DB" \
  --set ON_ERROR_STOP=1 <<'SQL'
\getenv app_password APP_DB_PASSWORD
\getenv migrator_password MIGRATOR_DB_PASSWORD

CREATE ROLE renalrelay_app LOGIN PASSWORD :'app_password';
CREATE ROLE renalrelay_migrator LOGIN PASSWORD :'migrator_password';
GRANT CONNECT ON DATABASE renalrelay TO renalrelay_app, renalrelay_migrator;
GRANT USAGE ON SCHEMA public TO renalrelay_app, renalrelay_migrator;
GRANT CREATE ON SCHEMA public TO renalrelay_migrator;
ALTER DEFAULT PRIVILEGES FOR ROLE renalrelay_migrator IN SCHEMA public
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO renalrelay_app;
ALTER DEFAULT PRIVILEGES FOR ROLE renalrelay_migrator IN SCHEMA public
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO renalrelay_app;
SQL

unset APP_DB_PASSWORD MIGRATOR_DB_PASSWORD
