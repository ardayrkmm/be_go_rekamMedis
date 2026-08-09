#!/bin/bash
# =============================================================
# Script Deploy Backend Go (MySQL) ke VPS
# Jalankan di VPS: bash deploy.sh
# =============================================================

set -e

PROJECT_DIR="/home/ardas/be_go_rekamMedis"
SERVICE_NAME="backend_rekammedis"
DB_NAME="rekam_medis"
REPO_URL="https://github.com/ardayrkmm/be_go_rekamMedis.git"  # <-- GANTI dengan URL repo Git kamu

echo "======================================"
echo " DEPLOY BACKEND REKAM MEDIS (MySQL)"
echo "======================================"

# ── 1. Stop service lama ──────────────────────────────────────
echo ""
echo "[1/7] Menghentikan service lama..."
sudo systemctl stop $SERVICE_NAME 2>/dev/null || echo "  Service tidak aktif, lanjut..."

# ── 2. Pull/Clone kode terbaru ────────────────────────────────
echo ""
echo "[2/7] Update kode dari Git..."
if [ -d "$PROJECT_DIR/.git" ]; then
    cd $PROJECT_DIR
    git pull origin main
    echo "  Git pull selesai"
else
    echo "  Folder tidak ada, clone baru..."
    cd /home/ardas
    git clone $REPO_URL be_go_rekamMedis
    cd $PROJECT_DIR
fi

# ── 3. Setup database MySQL ───────────────────────────────────
echo ""
echo "[3/7] Setup database MySQL..."
mysql -u root -p<<'MYSQL_EOF'
CREATE DATABASE IF NOT EXISTS rekam_medis CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
MYSQL_EOF
echo "  Database rekam_medis siap"

# ── 4. Update .env untuk VPS ─────────────────────────────────
echo ""
echo "[4/7] Cek file .env..."
if [ ! -f "$PROJECT_DIR/.env" ]; then
    echo "  File .env tidak ada! Buat manual dulu:"
    echo "  nano $PROJECT_DIR/.env"
    echo ""
    echo "  Isi dengan:"
    echo "    APP_ENV=production"
    echo "    APP_PORT=8001"
    echo "    JWT_SECRET=ganti_dengan_secret_aman"
    echo "    DB_DSN=root:PASSWORD_MYSQL@tcp(127.0.0.1:3306)/rekam_medis?charset=utf8mb4&parseTime=True&loc=Local"
    exit 1
fi
echo "  .env ditemukan"

# ── 5. Build binary Linux ─────────────────────────────────────
echo ""
echo "[5/7] Build binary Go..."
cd $PROJECT_DIR

# Pastikan Go terinstall
if ! command -v go &> /dev/null; then
    echo "  Go belum terinstall! Install dulu:"
    echo "  wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz"
    echo "  sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz"
    echo "  export PATH=\$PATH:/usr/local/go/bin"
    exit 1
fi

go mod tidy
go build -o api_linux ./cmd/api
echo "  Build selesai: api_linux"

# ── 6. Chmod dan set permission ───────────────────────────────
echo ""
echo "[6/7] Set permission..."
chmod +x api_linux

# ── 7. Buat/Update systemd service ───────────────────────────
echo ""
echo "[7/7] Update systemd service..."
sudo tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null <<EOF
[Unit]
Description=Backend Rekam Medis Golang (MySQL)
After=network.target mysql.service

[Service]
Type=simple
User=ardas
WorkingDirectory=$PROJECT_DIR
ExecStart=$PROJECT_DIR/api_linux
EnvironmentFile=$PROJECT_DIR/.env
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable $SERVICE_NAME
sudo systemctl start $SERVICE_NAME

echo ""
echo "======================================"
echo " Deploy selesai!"
echo " Cek status: sudo systemctl status $SERVICE_NAME"
echo " Lihat log: sudo journalctl -u $SERVICE_NAME -f"
echo "======================================"
