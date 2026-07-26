# Get a copy of the go compiler
FROM golang:1.25-alpine AS builder

# Set working directory. "." now means this location
WORKDIR /src

# Copy over the source code. 
COPY ./main.go .

# Compile the program
RUN go build -o /bin/main ./main.go

# Start the run stage. 
# It might be possible to run inside of the fronted docker file instead
FROM scratch

# Copy the binary over
COPY --from=builder /bin/main /bin/main

#Allow connections to the api port
EXPOSE 8090

# Run it
CMD ["/bin/main"]
