# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# ติดตั้ง dependencies ก่อน (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary (static binary, ไม่ต้องพึ่ง C library)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/api ./cmd/main.go

# Stage 2: Runtime (image เล็กมาก ~5MB)
FROM scratch

# Copy SSL certificates (สำหรับ HTTPS calls)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /app/api /api

# Non-root user
USER 65534

EXPOSE 8080

ENTRYPOINT ["/api"]
