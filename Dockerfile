FROM golang:1.24-alpine AS builder
WORKDIR /app
ARG APP_NAME="ticket-reservation-service"
RUN apk add --no-cache git

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ${APP_NAME} .
# ---
FROM alpine:edge
WORKDIR /app
ARG APP_NAME="ticket-reservation-service"
RUN apk add --no-cache ca-certificates tzdata
RUN adduser -D -g '' appuser
COPY --from=builder /app/${APP_NAME} .
COPY --from=builder /app/db/migrations /app/db/migrations
USER appuser

CMD ["./ticket-reservation-service", "serve"]