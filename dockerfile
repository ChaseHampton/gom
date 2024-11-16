# Use the official PostgreSQL image as the base
FROM postgres:latest

# Environment variables for PostgreSQL
ENV POSTGRES_USER=${POSTGRES_USER:-admin}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-password}
ENV POSTGRES_DB=${POSTGRESS_DB:-gom_db}

# Copy SQL script to create tables into the container
COPY ./sql/gom_base_create.sql /docker-entrypoint-initdb.d/

# Expose the default PostgreSQL port
EXPOSE 5432

# Set up persistent storage for PostgreSQL data
VOLUME /var/lib/postgresql/data

RUN useradd -ms /bin/bash gom_user
USER gom_user
