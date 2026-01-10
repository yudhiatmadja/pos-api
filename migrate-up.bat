@echo off
setlocal EnableDelayedExpansion

echo.
echo ====================================
echo  🚀 Database Migration Tool
echo ====================================
echo.

REM Cek apakah PostgreSQL container running
docker ps | findstr postgres-pos >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo ❌ PostgreSQL container tidak running!
    echo.
    echo Jalankan: docker-compose up -d postgres
    pause
    exit /b 1
)

echo ✅ PostgreSQL container detected
echo.
echo 🔄 Menjalankan migrasi...
echo.

docker run --rm ^
-v "%cd%/db/migration:/migrations" ^
--network pos-api_default ^
migrate/migrate:v4.19.1 ^
-path=/migrations ^
-database "postgres://postgres:password@postgres-pos:5432/mypos_db?sslmode=disable" ^
-verbose up

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ====================================
    echo  ✅ Migrasi berhasil dijalankan!
    echo ====================================
) else (
    echo.
    echo ====================================
    echo  ❌ Migrasi gagal!
    echo ====================================
    echo.
    echo Troubleshooting:
    echo 1. Pastikan PostgreSQL running: docker-compose up -d postgres
    echo 2. Cek password: password
    echo 3. Cek network: pos-api_default
    echo 4. Cek container: postgres-pos
)

echo.
pause