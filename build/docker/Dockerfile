############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
LABEL authors="Andrew Strongin"

# Install git (required for fetching dependencies)
RUN apk update && apk add --no-cache git

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

############################
# STEP 2 build a small image
############################
FROM alpine:latest

WORKDIR /app

# Copy the binary and scripts directory from builder
COPY --from=builder /app/ .

ENV HOST=0.0.0.0
ENV PORT=8080

EXPOSE 8080

# Run the binary
CMD ["./main"]