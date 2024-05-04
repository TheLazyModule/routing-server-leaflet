#!/bin/sh


# Fetch the .env file from S3
aws s3 cp s3://routing-env/app.env /app/app.env

# Export environment variables
set -o allexport
source /app/app.env
set +o allexport

# Start your application
exec "$@"
