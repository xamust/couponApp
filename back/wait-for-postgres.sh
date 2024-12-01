#!/bin/sh
# wait-for-postgres.sh - only for docker-compose

set -e

host="$1"
shift

until PGPASSWORD="$POSTGRES_PASSWORD" psql -h "$host" "$POSTGRES_DB" -U "$POSTGRES_USER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"

>&2 echo "Applying migrations..."
/app/migrator up

exec "$@"
