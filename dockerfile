
# Dockerfile
FROM golang:1.21 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a smaller base image for the final image
FROM alpine:latest

# Add ca certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]