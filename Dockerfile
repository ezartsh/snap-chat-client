FROM golang:alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
# COPY go.mod go.sum ./
# RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . /app

# Build
RUN go build -o snap-chat

# Define the command to run
ENTRYPOINT ["./snap-chat"]
