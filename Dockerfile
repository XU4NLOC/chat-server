# ── Build stage ──────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chat-server .

# ── Run stage ────────────────────────────────────────────────
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/chat-server .
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./chat-server"]
