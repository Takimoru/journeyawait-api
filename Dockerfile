# Use golang base image with specified version
FROM golang:1.22.2-alpine3.19 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application source code to the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:3.19

# Set the working directory inside the container
WORKDIR /root

# Copy the compiled executable from the previous stage
COPY --from=build /app/main .


# Command to run the executable
CMD ["./main"]