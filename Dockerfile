FROM golang:1.24.4-alpine AS builder
WORKDIR /app
ARG APP_NAME="ticket-reservation-service"

# Install git for dependency retrieval
RUN apk add --no-cache git

# Build environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -trimpath -ldflags="-s -w" -o ${APP_NAME} .
# ---
FROM alpine:latest
WORKDIR /app
ARG APP_NAME="ticket-reservation-service"

# Install required packages
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user and group
RUN addgroup -S appuser \
 && adduser -S -G appuser -H -s /sbin/nologin appuser

# Copy artifacts and give ownership to non-root user
COPY --from=builder --chown=appuser:appuser /app/${APP_NAME} /app/${APP_NAME}
COPY --from=builder --chown=appuser:appuser /app/db/migrations /app/db/migrations

# Switch to non-root user
USER appuser

# Run the application
CMD ["./ticket-reservation-service", "serve"]