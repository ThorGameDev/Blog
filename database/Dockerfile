# Pull the image
FROM postgres:18.4-alpine3.24

# Initialize the database with the starting sql file
COPY ./blog.sql /docker-entrypoint-initdb.d/
