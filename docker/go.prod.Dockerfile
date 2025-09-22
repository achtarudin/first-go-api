FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build only if binary does not exist
RUN if [ ! -f "/app/bin/main" ]; then \
      mkdir -p /app/bin && \
      CGO_ENABLED=0 go build -o /app/bin/main . ; \
    else \
      echo "Binary already exists, skip build" ; \
    fi


FROM alpine:3.22

# Menetapkan direktori kerja
WORKDIR /app

# Menyalin HANYA file biner 'main' yang sudah dicompile dari tahap 'builder'
# COPY  ./bin/main .
COPY --from=builder /app/bin/main .


# Mengekspos port yang digunakan oleh aplikasi
EXPOSE 8080

# Perintah untuk menjalankan aplikasi saat container dimulai
# Langsung menjalankan file binernya, bukan via 'air'
CMD ["./main"]
