# Menggunakan base image resmi Golang versi 1.21 berbasis Alpine
# Alpine dipilih karena ukurannya yang sangat kecil
FROM golang:1.24-alpine

# Menginstal package yang dibutuhkan:
# - git: diperlukan untuk 'go get' atau 'go install' dari repo
# - build-base: berisi compiler C dan tools lain yang mungkin dibutuhkan dependensi
RUN apk add --no-cache git build-base

# Menetapkan direktori kerja di dalam container
WORKDIR /app

# Menyalin file manajemen dependensi terlebih dahulu untuk memanfaatkan cache Docker
COPY go.mod go.sum ./
# Mengunduh semua dependensi
RUN go mod download

# Menyalin sisa kode sumber aplikasi Anda
COPY . .

# Menginstal 'air', sebuah tool untuk live-reloading aplikasi Go
RUN CGO_ENABLED=0 go install github.com/air-verse/air@latest

# Mengekspos port 8080 yang akan digunakan oleh aplikasi Go kita
EXPOSE 8080

# Perintah default saat container dijalankan
# 'air' akan memonitor perubahan file dan otomatis me-restart server
CMD ["air"]
