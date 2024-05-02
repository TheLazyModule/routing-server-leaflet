#!/bin/sh

set -e
echo "Starting the migration process..."
. /app/app.env   # Changed 'source' to '.'
migrate -path=/app/migrations -database "$DATABASE_URL" -verbose up
echo "Migration completed."

echo "Pre-inserting data..."
cd /app/geojson-graph  && conda run -n myenv python main.py && cd /app
echo "Data pre-insertion completed."

echo "Starting the main application..."
exec "$@"
echo "Main application started."
