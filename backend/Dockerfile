# Get a copy of the go compiler
FROM golang:1.26-alpine AS builder

# Set working directory. "." now means this location
WORKDIR /src

# Create user
RUN addgroup -S user --gid 1000 && adduser -S user -G user -u 1000

# Create the uploads volume with proper permissions
RUN mkdir -p /uploads && chown -R user:user /uploads
RUN mkdir -p /workdir && chown -R user:user /workdir

# switch to user
USER user

# Copy over dependency list
COPY --chown=user:user go.mod go.sum ./

# Compile the dependencies, and use a cache
RUN --mount=type=cache,target=/go/pkg/mod,uid=1000,gid=1000 \
    go mod download

# Copy over the source code.
COPY --chown=user:user . .

# Compile the program
RUN --mount=type=cache,target=/go/pkg/mod,uid=1000,gid=1000 \
    --mount=type=cache,target=/home/user/.cache/go-build,uid=1000,gid=1000 \
    CGO_ENABLED=0 GOOS=linux go build -o /home/user/main ./cmd/api

# Start the run stage.
# It might be possible to run inside of the fronted docker file instead
FROM scratch

#Allow connections to the api port
EXPOSE 8090

# Copy the user over
COPY --from=builder /etc/passwd /etc/passwd

# Switch to the user
USER user

# Copy the binary over
COPY --chown=user:user --from=builder /home/user/main /home/user/main

# Copy the uploads folder over
COPY --chown=user:user --from=builder /uploads /uploads
COPY --chown=user:user --from=builder /workdir /workdir

# Run it
CMD ["/home/user/main"]
