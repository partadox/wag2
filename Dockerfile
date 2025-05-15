# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev git

WORKDIR /app

# Copy semua file
COPY . .

# Build dengan CGO_ENABLED=1
RUN CGO_ENABLED=1 go build -o app

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata sqlite

# Siapkan direktori aplikasi
WORKDIR /app

# Copy binary dan file-file yang diperlukan
COPY --from=builder /app/app .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
# COPY --from=builder /app/.env ./.env

# Buat direktori data WhatsApp
RUN mkdir -p /app/whatsapp-data

# Set izin yang benar
RUN chmod +x /app/app

# Expose port aplikasi
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV LISTEN_ADDR=:8080
ENV WHATSAPP_DATA_DIR=/app/whatsapp-data

# Jalankan aplikasi
CMD ["./app"]