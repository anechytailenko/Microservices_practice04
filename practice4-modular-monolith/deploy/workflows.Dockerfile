# ==========================================
# STAGE 1: Builder
# ==========================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG COMMIT_HASH="unknown"

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X main.CommitHash=${COMMIT_HASH}" \
    -o workflows-service ./cmd/workflows/main.go

# ==========================================
# STAGE 2: Release
# ==========================================
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/workflows-service .

EXPOSE 8080

ENTRYPOINT ["./workflows-service"]