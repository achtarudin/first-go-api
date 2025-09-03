# first-go-api

### Cara Menggunakannya

1.  Simpan kedua file di atas (`Dockerfile` dan `docker-compose.yml`) di *root* direktori proyek Anda.
2.  Pastikan Docker sudah berjalan.

#### Untuk Menjalankan Lingkungan **Development**:
Buka terminal dan jalankan:
```bash
docker-compose up dev
```
Ini akan membuat kontainer dari image `golang:1.24-alpine`, me-mount kode Anda, dan menjalankan `air`. Buka `http://localhost:8081`. Setiap kali Anda menyimpan perubahan pada file `.go`, server akan otomatis restart.

#### Untuk Menjalankan Lingkungan **Production**:
Buka terminal dan jalankan:
```bash
docker-compose up --build apps