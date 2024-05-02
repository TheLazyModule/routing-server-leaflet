#!/bin/sh

set -e
echo "Starting the main application..."
exec "$@"
echo "Main application started."
