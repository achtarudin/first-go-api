# =================================================================
# TAHAP 1: BUILDER
# Tahap ini menggunakan image Go lengkap untuk mengkompilasi aplikasi.
# =================================================================
FROM golang:1.24-alpine AS builder

# Menetapkan direktori kerja
WORKDIR /app

# Menyalin file dependensi dan mengunduhnya untuk memanfaatkan cache
COPY go.mod go.sum ./
RUN go mod download

# Menyalin seluruh source code
COPY . .

# Mengkompilasi aplikasi Go.
# CGO_ENABLED=0 membuat biner yang statis (tidak bergantung pada library C sistem).
# -o /app/main menentukan outputnya adalah satu file bernama 'main'.
RUN CGO_ENABLED=0 go build -o /app/main .


# =================================================================
# TAHAP 2: FINAL
# Tahap ini menggunakan base image yang sangat kecil dan hanya berisi
# hasil kompilasi dari tahap sebelumnya.
# =================================================================
FROM alpine:latest

# Menetapkan direktori kerja
WORKDIR /app

# Menyalin HANYA file biner 'main' yang sudah dicompile dari tahap 'builder'
COPY --from=builder /app/main .

# Mengekspos port yang digunakan oleh aplikasi
EXPOSE 8080

# Perintah untuk menjalankan aplikasi saat container dimulai
# Langsung menjalankan file binernya, bukan via 'air'
CMD ["./main"]
