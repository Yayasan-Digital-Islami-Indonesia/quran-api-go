#!/bin/sh
set -e
/app/migrate
exec /app/api
