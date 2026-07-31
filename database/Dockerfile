# Pull the image
FROM postgres:18.4-alpine3.24

# Initialize the database with the starting sql file
COPY ./blog.sql /docker-entrypoint-initdb.d/

HEALTHCHECK --interval=1m --start-period=1m --start-interval=1s --timeout=5s --retries=3 \
    CMD pg_isready -U postgres || exit 1
