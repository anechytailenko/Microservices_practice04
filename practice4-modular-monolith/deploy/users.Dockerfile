ARG COMMIT_HASH="unknown"

RUN go build -ldflags="-X main.CommitHash=${COMMIT_HASH}" -o /app/server ./cmd/users

