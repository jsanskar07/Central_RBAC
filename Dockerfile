# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for downloading dependencies
RUN apk add --no-cache git

# Copy go mod and sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application statically (CGO_ENABLED=0 ensures it runs on any Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final Stage (Minimal Image)
FROM alpine:latest

WORKDIR /app

# Install CA certificates (important if you ever make outbound HTTPS requests)
RUN apk --no-cache add ca-certificates tzdata

# Copy the pre-built binary file from the builder stage
COPY --from=builder /app/main .

EXPOSE 9000

# Command to run the executable
CMD ["./main"]
