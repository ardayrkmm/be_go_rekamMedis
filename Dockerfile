# Tahap 1: Membangun Aplikasi (Builder Stage)
FROM golang:1.25-alpine AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Tambahkan dependensi sistem yang mungkin dibutuhkan (git, gcc, dll jika menggunakan CGO)
# Karena aplikasi Go biasanya tidak butuh banyak sistem dependensi jika CGO_ENABLED=0
RUN apk add --no-cache git

# Salin go.mod dan go.sum terlebih dahulu untuk memanfaatkan cache layer Docker
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode sumber (pastikan file sensitif masuk .dockerignore)
COPY . .

# Bangun file binary Go (CGO_ENABLED=0 membuat binary static yang berjalan lancar di alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Tahap 2: Menjalankan Aplikasi (Runner Stage)
FROM alpine:latest

WORKDIR /app

# Tambahkan tzdata untuk pengaturan timezone jika dibutuhkan
RUN apk add --no-cache tzdata

# Salin binary dari tahap builder
COPY --from=builder /app/main .

# Expose port (default 8001 sesuai setting di go)
EXPOSE 8001

# Jalankan aplikasi
CMD ["./main"]
