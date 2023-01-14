#!/bin/sh

set -e
#running in golang main file
#echo "run db migration"
#
#source /app/app.env
#
#/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up


echo "start the app"

exec "$@"

