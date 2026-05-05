FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy the module files and source code
COPY go.mod ./
COPY . .

# Run go mod tidy to resolve dependencies
RUN go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
