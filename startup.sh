#!/bin/sh

# Install AWS CLI if not already installed
if ! command -v aws &> /dev/null
then
    echo "AWS CLI could not be found, installing..."
    apt-get update && apt-get install -y awscli
fi

# Fetch the .env file from S3
aws s3 cp s3://routing-env/app.env /app/app.env

# Export environment variables
set -o allexport
source /your/app/app.env
set +o allexport

# Start your application
exec "$@"
