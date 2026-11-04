#!/bin/sh
set -eu

project_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
secret_dir="$project_root/secrets"

if [ -e "$secret_dir" ]; then
  printf '%s\n' 'secrets/ already exists; refusing to overwrite local credentials.' >&2
  exit 1
fi

if ! command -v openssl >/dev/null 2>&1; then
  printf '%s\n' 'OpenSSL is required to generate local development credentials.' >&2
  exit 1
fi

umask 077
mkdir -m 700 "$secret_dir"
trap 'rm -rf "$secret_dir"' EXIT HUP INT TERM
for secret_name in postgres_admin_password app_db_password migrator_db_password; do
  openssl rand -hex 32 > "$secret_dir/$secret_name"
  chmod 600 "$secret_dir/$secret_name"
done
if [ ! -f "$project_root/.env" ]; then
  sed \
    -e "s/^LOCAL_UID=.*/LOCAL_UID=$(id -u)/" \
    -e "s/^LOCAL_GID=.*/LOCAL_GID=$(id -g)/" \
    "$project_root/.env.example" > "$project_root/.env"
  chmod 600 "$project_root/.env"
fi

trap - EXIT HUP INT TERM

printf '%s\n' 'Generated local-only secrets with restrictive file permissions.'
