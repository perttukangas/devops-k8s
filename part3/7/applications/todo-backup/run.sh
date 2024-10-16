#!/bin/bash

if [[ -z "$POSTGRES_DB" || -z "$POSTGRES_USER" || -z "$POSTGRES_HOST" || -z "$POSTGRES_PASSWORD" || -z "$POSTGRES_PORT" || -z "$GCLOUD_TOKEN" ]]; then
  echo "One or more required environment variables are not set."
  echo "Please set POSTGRES_DB, POSTGRES_USER, POSTGRES_HOST, POSTGRES_PASSWORD, POSTGRES_PORT, and GCLOUD_TOKEN."
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

OBJECT_LOCATION="./todo_backup_$TIMESTAMP.sql"
OBJECT_CONTENT_TYPE="application/sql"
OBJECT_NAME="todo_backup_$TIMESTAMP.sql"
BUCKET_NAME="dwk-gke-bucket"

curl -X POST --data-binary @"$OBJECT_LOCATION" \
  -H "Authorization: Bearer $GCLOUD_TOKEN" \
  -H "Content-Type: $OBJECT_CONTENT_TYPE" \
  "https://storage.googleapis.com/upload/storage/v1/b/$BUCKET_NAME/o?uploadType=media&name=$OBJECT_NAME"