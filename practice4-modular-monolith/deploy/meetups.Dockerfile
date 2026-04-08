# Stage 1 : build

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /meetups-api ./cmd/api


# Stage 2: release

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /meetups-api .

EXPOSE 8080

CMD ["./meetups-api"]