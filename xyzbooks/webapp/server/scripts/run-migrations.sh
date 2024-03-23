#!/bin/sh

# Wait for the database to be ready
dockerize -wait tcp://${POSTGRES_HOST}:${POSTGRES_PORT} -timeout 30s

# Run migrations
echo "y" | /go/bin/migrate -source file://files -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable ${MIGRATE_COMMAND}