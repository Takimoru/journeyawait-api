# Menggunakan golang image resmi untuk build aplikasi
FROM golang:1.18 AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi
RUN go build -o main .

# Menyiapkan image yang lebih kecil untuk menjalankan aplikasi
FROM alpine:latest

# Install dependencies yang dibutuhkan
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy executable dari tahap build sebelumnya
COPY --from=builder /app/main .

# Menjalankan aplikasi
CMD ["./main"]
