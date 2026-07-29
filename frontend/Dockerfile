# Get a copy of node
FROM node:26-alpine AS builder

# Set directory. Only exists inside of docker container
WORKDIR /app

# Give the default user ownership of the /app directory
RUN chown node:node /app

#Do not run as root
USER node

# Copy the package.json file into the docker container /app directory
COPY package*.json ./

# Install the dependencies. Will run if package.json changes
RUN --mount=type=cache,target=/home/node/.npm,uid=1000,gid=1000 \
    npm ci --perfer-online

# Copy the directories that will be compiled/minified to the development environment
COPY ./scripts ./scripts
COPY ./html ./html
COPY ./style ./style

# Set the docker file to production mode by default. Allow CLI aruguments to change
ARG MODE=prod

# Run the bundling scripts. Should create a ./app/dist folder
RUN if [ "$MODE" = "dev" ]; then npm run build:dev; else npm run build:prod; fi


# Run stage
#FROM nginx-unprivileged:1.31-alpine3.24-perl
FROM nginxinc/nginx-unprivileged:1.31-alpine3.24-perl

# Copy the nginx config file over
COPY ./nginx.conf /etc/nginx/nginx.conf

# Copy the base directory. Should contain things like images or other resources
COPY ./root /usr/share/nginx/html
# Copy the compiled scripts to the /script directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Allow connections to the http protocol
EXPOSE 8080

# Launch nginx
CMD ["nginx", "-g", "daemon off;"]
