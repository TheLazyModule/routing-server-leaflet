# Golang build stage
FROM golang:alpine3.19 AS builder
# Install necessary tools for building and fetching dependencies
RUN apk add --no-cache curl tar
# Download and extract migrate tool, necessary for database migrations
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz -o migrate.tar.gz
RUN tar -xvzf migrate.tar.gz -C /usr/local/bin
WORKDIR /app
COPY . .
RUN go build -o main .

# Final stage
FROM continuumio/miniconda3:23.5.2-0-alpine

WORKDIR /app
# Copy built Go application and migrate tool from the builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
# Copy database migrations and environment configurations
COPY db/migrations /app/migrations
COPY app.env /app/
COPY .env /app/
COPY geojson-graph /app/geojson-graph
COPY public /app/public
COPY wait-for.sh startup.sh /app/
RUN chmod +x /app/wait-for.sh /app/startup.sh
RUN chmod +x /app/main
# Install netcat
RUN apk update && apk add --no-cache netcat-openbsd
# Copy the Conda environment file
COPY environment.yml /app/
# Create the Conda environment
RUN conda env create -f /app/environment.yml

# Activate the environment and install additional requirements if any
SHELL ["conda", "run", "-n", "myenv", "/bin/bash", "-c"]

# Continue with any further steps or modifications
EXPOSE 8080
ENTRYPOINT ["/app/startup.sh"]
CMD ["/app/main"]
