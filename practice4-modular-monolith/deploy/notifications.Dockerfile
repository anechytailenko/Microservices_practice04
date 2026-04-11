# ==========================================
# STAGE 1
# ==========================================

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG COMMIT_HASH="unknown"

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X main.CommitHash=${COMMIT_HASH}" \
    -o notification-service ./cmd/notifications/main.go

# ==========================================
# STAGE 2
# ==========================================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notification-service .

EXPOSE 8080

CMD ["./notification-service"]