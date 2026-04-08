# ==========================================
# STAGE 1: Builder
# ==========================================
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG COMMIT_HASH="unknown"

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X main.CommitHash=${COMMIT_HASH}" \
    -o users-service ./cmd/users/main.go

# ==========================================
# STAGE 2: Final Image
# ==========================================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/users-service .

EXPOSE 8080

CMD ["./users-service"]