#!/bin/bash

set -e
set -u

function create_user_and_database() {
  local database=$1
  echo "  Creating user and database '$database'"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
	    CREATE USER $database;
	    CREATE DATABASE $database;
	    GRANT ALL PRIVILEGES ON DATABASE $database TO $database;
EOSQL
  if [ -f "/docker-entrypoint-initdb.d/$database.sql" ]; then
    psql --username="$POSTGRES_USER" --dbname="$database" <"/docker-entrypoint-initdb.d/$database.sql"
  fi
}

if [ -n "$POSTGRES_DATABASES" ]; then
  echo "Multiple database creation requested: $POSTGRES_DATABASES"
  for db in $(echo "$POSTGRES_DATABASES" | tr ',' ' '); do
    create_user_and_database "$db"
  done
  echo "Multiple databases created"
fi
