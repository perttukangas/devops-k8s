#!/bin/bash

if [[ -z "$POSTGRES_DB" || -z "$POSTGRES_USER" || -z "$POSTGRES_HOST" || -z "$POSTGRES_PASSWORD" || -z "$POSTGRES_PORT" ]]; then
  echo "One or more required environment variables are not set."
  echo "Please set POSTGRES_DB, POSTGRES_USER, POSTGRES_HOST, POSTGRES_PASSWORD, and POSTGRES_PORT."
  exit 1
fi

TIMESTAMP=$(date +"%Y%m%d%H%M%S")
export PGPASSWORD="$POSTGRES_PASSWORD"
pg_dump -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" > "./todo_backup_$TIMESTAMP.sql"

if [[ $? -eq 0 ]]; then
  echo "Database backup successful."
else
  echo "Database backup failed."
  exit 1
fi